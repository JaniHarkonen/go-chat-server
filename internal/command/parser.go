package command

import (
	"errors"
	"fmt"
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
func Parse(str string) (*command, error) {
	commName, argString, found := strings.Cut(str, " ")

	if !found {
		return &command{
			name:      commName,
			arguments: []*argument{},
		}, nil
	}

	arguments := make([]*argument, 0)
	cursor := -1

	// Appends a new argument
	appendArgument := func(argType int) {
		argument := argumentBuilder.String()
		arguments = append(arguments, newArgument(argType, &argument))
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
					appendArgument(typeString)
					goto ValidString // Skip fail
				}

				argumentBuilder.WriteByte(stringChar)
				advance(1)

				stringChar, _ = lookAhead(cursor)
			}

			// String was opened, but not closed -> invalid
			return nil, newParserError("Encountered a non-closing string starting at position %p!", startCursor)

		ValidString:
			advance(1)
		} else if isNumber, isDecimalFound := isNumberChar(startChar); isNumber { // Handle numeric argument
			for number, _ := lookAhead(cursor); number != 0; number, _ = lookAhead(cursor) {
				if isNumber, isDecimal := isNumberChar(number); isNumber {
					if isDecimal {
						if isDecimalFound { // Decimal point was already found
							return nil, newParserError("Encountered an invalid number starting at position %p!", startCursor)
						}

						isDecimalFound = true
					}

					argumentBuilder.WriteByte(number)
					advance(1)
				} else {
					break
				}
			}

			appendArgument(typeNumber)
		} else {
			for char, _ := lookAhead(cursor); cursor != 0 && cursor != ' '; char, _ = lookAhead(cursor) {
				argumentBuilder.WriteByte(char)
				advance(1)
			}

			argument := argumentBuilder.String()

			// Boolean argument
			if argument == "true" || argument == "false" {
				appendArgument(typeBool)
			} else if argument == "null" { // Null argument
				appendArgument(typeNull)
			} else { // Type to be determined
				appendArgument(typeAmbiguous)
			}
		}

		// End of arguments
		if next, pos := lookAhead(cursor); next == 0 {
			break
		} else if next != ' ' { // Next character has to be the argument separator
			return nil, newParserError("Argument separator (whitespace) expected after argument/command at position %p!", pos)
		}

		advance(1)
	}

	fmt.Println("returning")
	return &command{
		name:      commName,
		arguments: arguments,
	}, nil
}
