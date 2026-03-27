package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	ircbot "github.com/recanman/irc-logbot/internal"
	"github.com/recanman/irc-logbot/packages/client"
)

func parseTime(timeFmt *string) (string, error) {
	if (*timeFmt)[0:5] == "time." {
		switch strings.ToLower((*timeFmt)[5:]) {
		case "layout":
			return time.Layout, nil
		case "ansic":
			return time.ANSIC, nil
		case "unixdate":
			return time.UnixDate, nil
		case "rubydate":
			return time.RubyDate, nil
		case "rfc822":
			return time.RFC822, nil
		case "rfc822z":
			return time.RFC822Z, nil
		case "rfc850":
			return time.RFC850, nil
		case "rfc1123":
			return time.RFC1123, nil
		case "rfc1123z":
			return time.RFC1123Z, nil
		case "rfc3339":
			return time.RFC3339, nil
		case "rfc3339nano":
			return time.RFC3339Nano, nil
		case "kitchen":
			return time.Kitchen, nil
		case "stamp":
			return time.Stamp, nil
		case "stampmilli":
			return time.StampMilli, nil
		case "stampmicro":
			return time.StampMicro, nil
		case "stampnano":
			return time.StampNano, nil
		case "datetime":
			return time.DateTime, nil
		case "dateonly":
			return time.DateOnly, nil
		case "timeonly":
			return time.TimeOnly, nil
		default:
			return "", errors.New(fmt.Sprintf(
				"time format '%s' unrecognized",
				*timeFmt))
		}
	}
	return *timeFmt, nil
}

func main() {
	// Define flags without default values
	serverPtr := flag.String("server", "", "IRC server address")
	portPtr := flag.Int("port", 0, "Port number")
	nicknamePtr := flag.String("nickname", "", "Nickname for the bot")
	channelPtr := flag.String("channels", "", "Comma-separated list of channels to join")
	fileNamePtr := flag.String("file", "log", "File name prefix for logging")
	sslPtr := flag.Bool("ssl", false, "Use SSL")
	allowInsecurePtr := flag.Bool("allow-insecure", false, "Allow insecure SSL")
	timeFormatPtr := flag.String("time-format", "time.Kitchen", "Format of the time in logs (takes a go time constant or a go time.Format)")

	// Parse flags
	flag.Parse()

	// Check if required flags are set
	if *serverPtr == "" || *portPtr == 0 || *nicknamePtr == "" || *channelPtr == "" || *fileNamePtr == "" {
		fmt.Println("Missing required flags. Please specify --server, --port, --nickname, --channels, and --file.")
		return
	}

	// Retrieve flag values
	server := *serverPtr
	port := *portPtr
	nickname := *nicknamePtr
	channelsStr := *channelPtr                  // Channels string separated by comma
	channels := strings.Split(channelsStr, ",") // Splitting the string into a slice
	ssl := *sslPtr
	allowInsecure := *allowInsecurePtr
	timeFormat, err := parseTime(timeFormatPtr)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create client with specified options
	client := client.Create(server, port, nickname, client.ClientOptions{
		Channels: channels,
	}, ssl, allowInsecure)

	fmt.Println("Connecting to server...")
	err = client.Connect()
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	ircbot.FromClient(client, *fileNamePtr, timeFormat)
	fmt.Println("Bot is running...")

	select {}
}
