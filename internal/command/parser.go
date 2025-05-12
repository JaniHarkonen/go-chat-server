package command

import (
	"errors"
	"strconv"
	"strings"
)

func newParserError(message string, cursor int) error {
	return errors.New(strings.ReplaceAll(message, "%p", strconv.Itoa(cursor)))
}

func isStringChar(char byte) (bool, byte) {
	return char == '"' || char == '\'' || char == '`', char
}

func isNumberChar(char byte) (bool, bool) {
	return (char >= '0' && char <= '9') || char == '.', char == '.'
}

var argumentBuilder strings.Builder

// Parses a given string into commands (command-argument combinations) by splitting
// the string along whitespaces and treating the first split as the command name and
// the rest of the splits as arguments to the command.
func Parse(str string) (*Command, error) {
	commName, argString, found := strings.Cut(str, " ")

	if !found {
		return &Command{
			name:      commName,
			arguments: []*Argument{},
		}, nil
	}

	arguments := make([]*Argument, 0)
	cursor := -1

	// Appends a new argument
	appendArgument := func(argType int, value string) {
		arguments = append(arguments, newArgument(argType, &value))
	}

	// Peek at next token after the given position
	lookAhead := func(position int) (byte, int) {
		fixedPosition := position + 1

		if fixedPosition < len(argString) {
			return argString[fixedPosition], fixedPosition
		}

		return 0, fixedPosition
	}

	// Advance the cursor by a given amount
	advance := func(amount int) int {
		cursor += amount
		return cursor
	}

	for startChar, startCursor := lookAhead(cursor); ; startChar, startCursor = lookAhead(cursor) {
		argumentBuilder.Reset()

		// Handle string argument
		if isString, _ := isStringChar(startChar); isString {
			advance(1)
			stringChar, _ := lookAhead(cursor)

			for stringChar != 0 {
				// String closed -> valid
				if stringChar == startChar {
					appendArgument(TypeString, argumentBuilder.String())
					goto ValidString // Skip fail
				}

				argumentBuilder.WriteByte(stringChar)
				advance(1)

				stringChar, _ = lookAhead(cursor)
			}

			// String was opened, but not closed -> invalid
			return nil, newParserError("encountered a non-closing string starting at position %p", startCursor)

		ValidString:
			advance(1)
		} else if isNumber, isDecimalFound := isNumberChar(startChar); isNumber { // Handle numeric argument
			for number, _ := lookAhead(cursor); number != 0; number, _ = lookAhead(cursor) {
				if isNumber, isDecimal := isNumberChar(number); isNumber {
					if isDecimal {
						if isDecimalFound { // Decimal point was already found
							return nil, newParserError("encountered an invalid number starting at position %p", startCursor)
						}

						isDecimalFound = true
					}

					argumentBuilder.WriteByte(number)
					advance(1)
				} else {
					break
				}
			}

			appendArgument(TypeNumber, argumentBuilder.String())
		} else {
			for char, _ := lookAhead(cursor); char != 0 && char != ' '; char, _ = lookAhead(cursor) {
				argumentBuilder.WriteByte(char)
				advance(1)
			}

			var argType int
			argValue := argumentBuilder.String()

			// Boolean argument
			if argValue == "true" || argValue == "false" {
				argType = TypeBool
			} else if argValue == "null" { // Null argument
				argType = TypeNull
			} else if argValue[0] == '@' { // User identifier
				argType = TypeUser
				argValue = argValue[1:]
			} else { // Type to be determined
				argType = TypeAmbiguous
			}

			appendArgument(argType, argValue)
		}

		// End of arguments
		if next, pos := lookAhead(cursor); next == 0 {
			break
		} else if next != ' ' { // Next character has to be the argument separator
			return nil, newParserError("argument separator (whitespace) expected after argument/command at position %p", pos)
		}

		advance(1)
	}

	return &Command{
		name:      commName,
		arguments: arguments,
	}, nil
}
