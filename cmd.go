package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

type Args []string
type CmdFunc func(args Args)
type Command struct {
	Labels []string
	Func   CmdFunc
	Desc   string
}

var (
	cmds []Command
)

func init() {
	cmds = []Command{
		{[]string{"help", "h", "-h", "--help"}, HelpFunc, "Print instructions of how to to use passs"},
		{[]string{"list", "l"}, ListFunc, "List saved passwords"},
		{[]string{"insert", "i"}, InsertFunc, "Enter a new password to save"},
		{[]string{"retrieve", "r"}, RetrieveFunc, "Retrieve a password"},
		{[]string{"modify", "m"}, ModifyFunc, "Modify a record"},
		{[]string{"remove", "rm"}, RemoveFunc, "Remove a record from storage"},
	}
}

func FindRecord() int {
	var website string
	fmt.Print("Enter website to search for: ")
	fmt.Scanln(&website)

	var (
		wss   []string
		recis []int
	)

RecordsLoop:
	for i, rec := range Records {
		if !strings.Contains(rec.Website, website) {
			continue
		}

		recis = append(recis, i)

		for _, ws := range wss {
			if ws == rec.Website {
				continue RecordsLoop
			}
		}

		wss = append(wss, rec.Website)
	}

	if len(wss) == 0 {
		fmt.Println("No such website in storage!")
		os.Exit(1)
	}

	if len(wss) > 1 {
		fmt.Println("Found multiple websites:")
		fmt.Println()
		fmt.Print("\t")
		for _, ws := range wss {
			fmt.Printf("%v  ", ws)
		}
		fmt.Println()
		fmt.Println()
		fmt.Println("Please, clearify your query")
		return FindRecord()
	}

	var username string
	if len(recis) > 1 {
		fmt.Println("Found multiple accounts: ")
		fmt.Println()
		for _, i := range recis {
			Records[i].PrintHead()
		}
		fmt.Println()
		fmt.Print("Please, enter the needed username: ")
		fmt.Scanln(&username)
	}

	fmt.Println()
	for _, i := range recis {
		if strings.Contains(Records[i].Username, username) {
			return i
		}
	}

	fmt.Println("No such record!")
	os.Exit(1)
	return 0
}

func HelpFunc(args Args) {
	fmt.Println("Usage: ")
	fmt.Println()
	fmt.Printf("\t%v <command>\n", os.Args[0])
	fmt.Println()

	fmt.Println("Available commands:")
	fmt.Println()
	for _, cmd := range cmds {
		fmt.Printf("\t%-14v %v\n", cmd.Labels[0], cmd.Desc)
	}
	fmt.Println()
}

func ListFunc(args Args) {
	PrepareFile()

	fmt.Println("Records in storage: ")
	fmt.Println()
	for _, rec := range Records {
		rec.PrintHead()
	}
	fmt.Println()
}

func promt(s string) {
	fmt.Printf("%19v", s)
}

func InsertFunc(args Args) {
	PrepareFile()

	fmt.Println("Inserting a new password record")
	rec := PasssRecord{}

	promt("Website: ")
	fmt.Scanln(&rec.Website)

	promt("Username: ")
	fmt.Scanln(&rec.Username)

	for {
		promt("Password: ")
		pwd := string(gopass.GetPasswd())

		promt("Repeat: ")
		pwd2 := string(gopass.GetPasswd())

		if pwd == pwd2 {
			rec.Password = pwd
			break
		} else {
			fmt.Println("Entered passwords are different! Please try again.")
		}
	}

	fmt.Println("Inserted a new record:")
	rec.PrintHead()

	Records = append(Records, rec)
	EncryptAndSave()
}

func RetrieveFunc(args Args) {
	PrepareFile()

	fmt.Println("Retrieving a password from storage")
	rec := &Records[FindRecord()]
	rec.LoadPassword()

	fmt.Println("Password for: ")
	fmt.Println()
	rec.PrintHead()
	fmt.Println()
	fmt.Println("Has been put into your clipboard")
}

func ModifyFunc(args Args) {
	PrepareFile()

	fmt.Println("Modifying a record")
	rec := &Records[FindRecord()]

	fmt.Println()
	rec.PrintHead()
	fmt.Println()

	var (
		tmp  string
		tmp2 string
	)

	fmt.Println("Enter new values for field (or leave blank to keep them unchanged)")

	promt("Website: ")
	if _, err := fmt.Scanln(&tmp); err == nil {
		rec.Website = tmp
	} else {
		fmt.Println("Website field is unchanged")
	}

	promt("Username: ")
	if _, err := fmt.Scanln(&tmp); err == nil {
		rec.Username = tmp
	} else {
		fmt.Println("Username field is unchanged")
	}

	for {
		promt("Password: ")
		if _, err := fmt.Scanln(&tmp); err != nil {
			fmt.Println("Password is unchanged")
			break
		}

		promt("Repeat: ")
		fmt.Scanln(&tmp2)
		if tmp == tmp2 {
			rec.Password = tmp
			break
		} else {
			fmt.Println("Passwords are mismatched!")
		}
	}

	fmt.Println("Record has been successfully modified")
	EncryptAndSave()
}

func RemoveFunc(args Args) {
	PrepareFile()

	fmt.Println("Removing a record")
	index := FindRecord()
	rmrec := &Records[index]

	fmt.Println("Are you sure you want to delete the record?")
	fmt.Println()
	rmrec.PrintHead()
	fmt.Println()
	fmt.Print("[Y/n] ")

	var ans string
	fmt.Scanln(&ans)

	if ans != "Y" && ans != "y" {
		fmt.Println("Operation canceled")
		return
	}

	Records = append(Records[:index], Records[index+1:]...)
	fmt.Println("Record successfully removed!")
	EncryptAndSave()
}
