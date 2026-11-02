package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rotscher/httploom/loomclient/internal/api"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "loomclient:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("loomclient", flag.ContinueOnError)
	server := fs.String("server", envOr("LOOM_SERVER", "http://localhost:8080"), "loomserver base URL")
	timeout := fs.Duration("timeout", 10*time.Second, "request timeout")
	fs.Usage = func() { usage(fs) }
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		fs.Usage()
		return errors.New("missing command")
	}

	client, err := api.NewClientWithResponses(*server)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	switch cmd := fs.Arg(0); cmd {
	case "health":
		return health(ctx, client)
	case "weave":
		return create(ctx, client, fs.Args()[1:])
	case "list":
		return list(ctx, client)
	case "get":
		if fs.NArg() < 2 {
			return errors.New("get requires a weave id")
		}
		return get(ctx, client, fs.Arg(1))
	default:
		fs.Usage()
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func health(ctx context.Context, client *api.ClientWithResponses) error {
	resp, err := client.GetHealthWithResponse(ctx)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return unexpected(resp.StatusCode(), resp.Body)
	}
	return printJSON(resp.JSON200)
}

func create(ctx context.Context, client *api.ClientWithResponses, args []string) error {
	fs := flag.NewFlagSet("weave", flag.ContinueOnError)
	width := fs.Int("width", 80, "line width")
	if err := fs.Parse(args); err != nil {
		return err
	}

	text := strings.Join(fs.Args(), " ")
	if text == "" {
		stdin, err := readStdin()
		if err != nil {
			return errors.New("provide text as arguments or on stdin")
		}
		text = stdin
	}
	if strings.TrimSpace(text) == "" {
		return errors.New("text must not be empty")
	}

	resp, err := client.CreateWeaveWithResponse(ctx, api.CreateWeaveJSONRequestBody{
		Text:  text,
		Width: width,
	})
	if err != nil {
		return err
	}
	if resp.JSON201 == nil {
		if resp.JSON400 != nil {
			return errors.New(resp.JSON400.Message)
		}
		return unexpected(resp.StatusCode(), resp.Body)
	}

	fmt.Fprintf(os.Stderr, "created weave %s\n", resp.JSON201.Id)
	fmt.Println(resp.JSON201.Text)
	return nil
}

func list(ctx context.Context, client *api.ClientWithResponses) error {
	resp, err := client.ListWeavesWithResponse(ctx)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return unexpected(resp.StatusCode(), resp.Body)
	}
	return printJSON(resp.JSON200)
}

func get(ctx context.Context, client *api.ClientWithResponses, id string) error {
	resp, err := client.GetWeaveWithResponse(ctx, id)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		if resp.JSON404 != nil {
			return errors.New(resp.JSON404.Message)
		}
		return unexpected(resp.StatusCode(), resp.Body)
	}
	return printJSON(resp.JSON200)
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func readStdin() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", errors.New("stdin is a terminal")
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func unexpected(status int, body []byte) error {
	return fmt.Errorf("unexpected response %d: %s", status, strings.TrimSpace(string(body)))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func usage(fs *flag.FlagSet) {
	_, _ = fmt.Fprintf(fs.Output(), `Usage: loomclient [flags] <command> [args]

Commands:
  health            Check server health
  weave [-width N] <text...>   Weave text (or read it from stdin)
  list              List stored weaves
  get <id>          Fetch a single weave

Flags:
`)
	fs.PrintDefaults()
}
