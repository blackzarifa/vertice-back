package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			content, err := ioutil.ReadFile(filepath.Join("migrations", file.Name()))
			if err != nil {
				return err
			}

			if err := executeSQLFile(db, string(content), file.Name()); err != nil {
				return err
			}
			log.Printf("Executed migration: %s", file.Name())
		}
	}
	return nil
}

func executeSQLFile(db *sql.DB, content string, filename string) error {
	if strings.Contains(content, "DELIMITER") {
		return executeWithDelimiter(db, content, filename)
	}

	statements := strings.Split(content, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("error executing %s: %w", filename, err)
		}
	}
	return nil
}

func executeWithDelimiter(db *sql.DB, content string, filename string) error {
	delimiterRegex := regexp.MustCompile(`(?i)DELIMITER\s+(\S+)`)
	matches := delimiterRegex.FindAllStringSubmatch(content, -1)
	
	if len(matches) == 0 {
		return executeSQLFile(db, content, filename)
	}

	parts := delimiterRegex.Split(content, -1)
	delimiters := []string{";"}
	for _, match := range matches {
		delimiters = append(delimiters, match[1])
	}

	for i, part := range parts {
		if i >= len(delimiters) {
			break
		}
		
		delimiter := delimiters[i]
		if delimiter == "//" {
			statements := strings.Split(part, delimiter)
			for _, stmt := range statements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" || stmt == ";" {
					continue
				}
				if _, err := db.Exec(stmt); err != nil {
					return fmt.Errorf("error executing %s: %w", filename, err)
				}
			}
		} else {
			if err := executeSQLFile(db, part, filename); err != nil {
				return err
			}
		}
	}
	return nil
}
