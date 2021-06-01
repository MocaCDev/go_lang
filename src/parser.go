package src

import (
	"fmt"
	"log"
	//"encoding/json"
	"go_lang/go_lang_packager/packager"
)

type Settings struct {
	inter_access	bool
}

type Parser struct {
	lexer		interface{}
	lang_info	interface{}
	settings	Settings
}

func Init_parser(lexer interface{}) *Parser {
	return &Parser{ lexer: lexer, settings: Settings{inter_access: false} }
}

func (parser *Parser) get_next_token(lexer *Lexer) *Lexer {
	return lexer.Lex()
}

func (parser *Parser) parse_require(lexer *Lexer) *Parser {
	parser.get_next_token(lexer)


	if lexer.Current_Token == T_STR {

		info := packager.Read_info_package(lexer.Token_value)

		switch lexer.Token_value {
			case "lang_info.json": parser.lang_info = info.(*packager.LangInfo)
		}

		parser.get_next_token(lexer)

		if lexer.Current_Token == T_SEMI {
			parser.get_next_token(lexer)
		}

		parser.lexer = lexer

		return parser
	}

	log.Fatal("Parsing Error")
	return parser
}

// Function prone to change
func (parser *Parser) parse_print(lexer *Lexer) *Parser {
	parser.get_next_token(lexer)

	switch lexer.Current_Token {
		case K_VAR_N: {
			switch lexer.Token_value {
				case "inter": {
					if parser.settings.inter_access == false {
						log.Fatal("You do not have permission to use inter")
					}
					// Strict Syntax(For now)
					parser.get_next_token(lexer)

					switch lexer.Token_value {
						case "LangInfo": {
							parser.get_next_token(lexer)

							//parser.get_next_token(lexer)
							LI := parser.lang_info.(*packager.LangInfo)

							switch lexer.Token_value {
								case "lang_name":fmt.Println(LI.LangName)
								case "lang_version":fmt.Println(LI.LangVersion)
							}
						}
					}
				}
			}
		}
		case T_STR: {
			fmt.Println(lexer.Token_value)
		}
	}

	parser.get_next_token(lexer)

	return parser
}

func (parser *Parser) parse_wrap(lexer *Lexer) *Parser {
	parser.get_next_token(lexer)

	if lexer.Current_Token == T_LB {
		parser.get_next_token(lexer)

		switch lexer.Token_value {
			case "allow": {
				parser.get_next_token(lexer)

				switch lexer.Token_value {
					case "inter_access": {
						parser.settings.inter_access = true
						parser.get_next_token(lexer)
					}
				}

				if lexer.Current_Token == T_RB {
					parser.get_next_token(lexer)

					return parser
				}
				log.Fatal(fmt.Sprintf("Expected ']', got `%s` instead", lexer.Token_value))
			}
			log.Fatal(fmt.Sprintf("Expected ')', got `%s` instead", lexer.Token_value))
		}
		log.Fatal("Error parsing wrapper")
	}

	log.Fatal("Invalid syntax to wrapper")
	return parser;
}

func (parser *Parser) Parse() *Parser {

	lexer := parser.lexer.(*Lexer)
	for {
		switch lexer.Current_Token {
			case K_REQ: parser.parse_require(lexer)
			case K_PRINT: parser.parse_print(lexer)
			case T_WRAP: parser.parse_wrap(lexer)
			case T_EOF: return parser
			default: return parser
		}
	}

	return parser
}
