package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

func get_random_line(open_file *os.File) string {

	// Load it all into memory,
	// a better way would be to have the same byte count on each line and just get a multiple
	var lines []string
	scanner := bufio.NewScanner(open_file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	random_seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	line_number := random_seed.Intn(len(lines))

	//log.Printf("Picking line %s at line %d among total %d\n", lines[line_number], line_number, len(lines))

	return lines[line_number]
}

func define_client_id() string {
	// Get random line from a color wordlist
	color_file, color_err := os.Open("../wordlist/color.txt")
	// TODO check color_err
	if color_err != nil {
		return ""
	}
	defer color_file.Close()
	color_picked := get_random_line(color_file)
	// Get random line from an animal wordlist
	animal_file, animal_err := os.Open("../wordlist/animal.txt")
	// TODO check animal_err
	if animal_err != nil {
		return ""
	}
	defer animal_file.Close()
	animal_picked := get_random_line(animal_file)

	// return both string concat
	//log.Printf("Created new client id : %s_%s\n", color_picked, animal_picked)
	return color_picked + "_" + animal_picked
}
