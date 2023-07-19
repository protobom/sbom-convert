package readme

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	baseName := strings.Split(name, " ")[0]
	// baseName := split[0]
	cmdName := cmd.Name()
	if flags.HasAvailableFlags() {
		disc := fmt.Sprintf("Flags for `%s` subcommand", cmdName)
		if baseName == cmdName {
			disc = fmt.Sprintf("Flags for `%s`", baseName)
		}
		buf.WriteString(fmt.Sprintf("### Optional flags \n%s\n\n\n", disc))
		lines := FlagUsagesWrapped(flags, 0)
		lines = ReplacOutputDirHash(lines, baseName, "\" |")
		buf.WriteString(lines)
		buf.WriteString("\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		disc := fmt.Sprintf("Flags for all `%s` subcommands", baseName)
		buf.WriteString(fmt.Sprintf("### Global options flags\n%s\n\n\n", disc))
		lines := FlagUsagesWrapped(parentFlags, 0)
		lines = ReplacOutputDirHash(lines, baseName, "\" |")
		buf.WriteString(lines)
		buf.WriteString("\n\n")
	}
	return nil
}

// Hack to replace the host pushed to documentation
// 2DO how do we do this for cli command usage readme as well?
// m1 := regexp.MustCompile(`/home/.*`)
func ReplacOutputDirHash(src, appName, suffix string) string {
	m1 := regexp.MustCompile(`/home/.*`)
	readctedConfigStr := src
	replaceString := fmt.Sprintf("$${XDG_CACHE_HOME}/%s%s", appName, suffix)
	return m1.ReplaceAllString(readctedConfigStr, replaceString)
}

// FlagUsagesWrapped returns a string containing the usage information
// for all flags in the FlagSet. Wrapped to `cols` columns (0 for no
// wrapping)
func FlagUsagesWrapped(f *pflag.FlagSet, cols int) string {
	buf := new(bytes.Buffer)

	lines := make([]string, 0, 100)

	// maxlen := 0

	lines = append(lines, "| Short | Long | Description | Default |")
	lines = append(lines, "| --- | --- | --- | --- |")

	f.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		line := ""
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			line = fmt.Sprintf("| -%s | --%s |", flag.Shorthand, flag.Name)
		} else {
			line = fmt.Sprintf("| | --%s |", flag.Name)
		}

		_, usage := pflag.UnquoteUsage(flag)
		// if varname != "" {
		// 	line += fmt.Sprintf(" %s |", varname)
		// }
		if flag.NoOptDefVal != "" {
			switch flag.Value.Type() {
			case "string":
				line += fmt.Sprintf("  \"%s\" |", flag.NoOptDefVal)
			case "bool":
				if flag.NoOptDefVal != "true" {
					line += fmt.Sprintf(" %s |", flag.NoOptDefVal)
				}
			case "count":
				if flag.NoOptDefVal != "+1" {
					line += fmt.Sprintf(" %s |", flag.NoOptDefVal)
				}
			default:
				line += fmt.Sprintf(" %s |", flag.NoOptDefVal)
			}
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		// line += "\x00"
		// if len(line) > maxlen {
		// 	maxlen = len(line)
		// }

		line += fmt.Sprintf(" %s |", usage)
		if !defaultIsZeroValue(flag) {
			if flag.Value.Type() == "string" {
				line += fmt.Sprintf(" %q |", flag.DefValue)
			} else {
				line += fmt.Sprintf(" %s |", flag.DefValue)
			}
		} else {
			line += fmt.Sprintf(" |")
		}
		// if len(flag.Deprecated) != 0 {
		// 	line += fmt.Sprintf(" (DEPRECATED: %s) |", flag.Deprecated)
		// }

		lines = append(lines, line)
	})

	for _, line := range lines {
		// fmt.Println("###### line", line)
		// sidx := strings.Index(line, "\x00")
		// spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		// fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
		fmt.Fprintln(buf, line)
	}

	return buf.String()
}

// defaultIsZeroValue returns true if the default value for this flag represents
// a zero value.
func defaultIsZeroValue(f *pflag.Flag) bool {

	switch f.Value.Type() {
	case "bool":
		return f.DefValue == "false"
	// case *durationValue:
	// 	// Beginning in Go 1.7, duration zero values are "0s"
	// 	return f.DefValue == "0" || f.DefValue == "0s"
	// case *intValue, *int8Value, *int32Value, *int64Value, *uintValue, *uint8Value, *uint16Value, *uint32Value, *uint64Value, *countValue, *float32Value, *float64Value:
	// 	return f.DefValue == "0"
	case "int", "count":
		return f.DefValue == "0"
	case "string":
		return f.DefValue == ""
	// case *ipValue, *ipMaskValue, *ipNetValue:
	// 	return f.DefValue == "<nil>"
	case "intSlice", "stringSlice":
		return f.DefValue == "[]"
	default:
		fmt.Println("#### unknown", f.Name, f.Value.Type())
		switch f.Value.String() {
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		}
		return false
	}

}

// Test to see if we have a reason to print See Also information in docs
// Basically this is a test for a parent command or a subcommand which is
// both not deprecated and not the autogenerated help command.
func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

// GenMarkdownCustom creates custom markdown output.
func GenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return err
	}

	if len(cmd.Example) > 0 {
		buf.WriteString(fmt.Sprintf("### Examples for running `%s`\n\n", name))
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if hasSeeAlso(cmd) {
		buf.WriteString("### SEE ALSO\n\n")
		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := pname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			buf.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", pname, linkHandler(link), parent.Short))
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			buf.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", cname, linkHandler(link), child.Short))
		}
		buf.WriteString("\n")
	}
	// if !cmd.DisableAutoGenTag {
	// 	buf.WriteString("###### Auto generated by spf13/cobra on " + time.Now().Format("2-Jan-2006") + "\n")
	// }
	_, err := buf.WriteTo(w)
	return err
}
