package trapcmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/bugout-dev/bugout-go/pkg/spire"
)

type TrapEntry struct {
	Title   string
	Content string
	Tags    []string
	Context spire.EntryContext
}

func Render(invocation []string, result InvocationResult, title string, tags []string, showEnv bool) TrapEntry {
	renderedTitle := RenderTitle(invocation, result, title, tags, showEnv)
	renderedTags := RenderTags(invocation, result, title, tags, showEnv)
	renderedContent := RenderContent(invocation, result, title, tags, showEnv)
	return TrapEntry{Title: renderedTitle, Content: renderedContent, Tags: renderedTags, Context: spire.EntryContext{ContextType: "trap"}}
}

func RenderTitle(invocation []string, result InvocationResult, title string, tags []string, showEnv bool) string {
	if title != "" {
		return title
	}
	return fmt.Sprintf("Command: %s (exited with code %d)", invocation[0], result.ExitCode)
}

func RenderTags(invocation []string, result InvocationResult, title string, tags []string, showEnv bool) []string {
	trapTags := []string{
		"trap",
		"cli",
		fmt.Sprintf("code:%d", result.ExitCode),
		fmt.Sprintf("exit:%d", result.ExitCode),
		fmt.Sprintf("env:%v", showEnv),
	}
	finalTags := append(trapTags, tags...)
	return finalTags
}

func RenderContent(invocation []string, result InvocationResult, title string, tags []string, showEnv bool) string {
	quotedInvocation := make([]string, len(invocation))
	for i, component := range invocation {
		quotedInvocation[i] = strconv.Quote(component)
	}

	envvars := os.Environ()
	sort.Strings(envvars)
	quotedEnvvars := make([]string, len(envvars))
	for i, envvar := range envvars {
		quotedEnvvars[i] = fmt.Sprintf("`%s`", envvar)
	}

	var content string = strings.Join([]string{
		fmt.Sprintf("## invocation\n```\n%s\n```\n", strings.Join(quotedInvocation, " ")),
		fmt.Sprintf("## exit code\n`%d`\n", result.ExitCode),
		fmt.Sprintf("## stdout\n```\n%s\n```\n", result.OutBuffer.String()),
		fmt.Sprintf("## stderr\n```\n%s\n```\n", result.ErrBuffer.String()),
	}, "\n")

	if showEnv {
		content = strings.Join([]string{content, fmt.Sprintf("## env\n- %s\n", strings.Join(quotedEnvvars, "\n- "))}, "\n")
	}

	return content
}
