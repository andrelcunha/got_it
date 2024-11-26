package init

import (
	"fmt"
	"got_it/internal/config"
	"os"
)

type Init struct {
	conf *config.Config
}

func NewInit() *Init {
	conf := config.NewConfig()
	return &Init{
		conf: conf,
	}
}

func (i *Init) IsInitialized() bool {
	if _, err := os.Stat(i.conf.GetGotDir()); os.IsNotExist(err) {
		fmt.Println("Not a Got_it repository. Run 'got init' first.")
		return false
	}
	return true
}

func (i *Init) InitRepo() {
	gotDir := i.conf.GetGotDir()

	// Check if the .got directory already exists
	if _, err := os.Stat(gotDir); !os.IsNotExist(err) {
		fmt.Println("Repository already initialized.")
		return
	}

	// Create the .got directory
	if err := os.Mkdir(gotDir, 0755); err != nil {
		fmt.Println("Error creating .got directory:", err)
		return
	}
	// Create the .got/objects directory
	if err := os.Mkdir(gotDir+"/objects", 0755); err != nil {
		fmt.Println("Error creating .got/objects directory:", err)
		return
	}

	fmt.Println("Repository initialized successfully.")

}
