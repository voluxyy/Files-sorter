package main

import (
	"C"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println(os.Executable())
	myApp := app.New()
	myWindow := myApp.NewWindow("Files sorter")

	iconLink, errIcon := fyne.LoadResourceFromURLString("https://github.com/voluxyy/Files-sorter/blob/main/images/icon.png?raw=true")
	if errIcon != nil {
		log.Println(errIcon)
	}

	myWindow.SetIcon(iconLink)
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.CenterOnScreen()
	myWindow.FixedSize()

	var repository string
	var character string

	label := widget.NewLabel("Répertoire :")
	label.Wrapping = fyne.TextWrapWord

	label1 := widget.NewLabel("Caractères :")

	label2 := widget.NewLabel(`Résultat : En attente des infos..`)

	label3 := widget.NewLabel("Gestion d'erreur : ")
	label3.Wrapping = fyne.TextWrapWord

	inputCharacter := widget.NewEntry()
	inputCharacter.SetPlaceHolder("Entrer les caractères...")

	var dirInfo fyne.ListableURI
	buttonFolder := widget.NewButton("Dossier", func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				label3.SetText(`Gestion d'erreur : Problème avec le selecteur de répertoire.`)
			}
			if dir.Authority() == "inutile pour le moment" {
				label3.SetText(`Gestion d'erreur : Problème d'autorité d'accès au répertoire.`)
			}
			dirInfo = dir
			label.SetText("Répetoire : " + dirInfo.Path())
		}, myWindow)
	})

	var checked bool
	check := widget.NewCheck("Créer le dossier sur le bureau", func(isCheck bool) {
		if isCheck == true {
			checked = true
		} else {
			checked = false
		}
	})

	launchButton := widget.NewButton("Lancer", func() {
		// Update label status
		label2.SetText(`Résultat : En cours d'exécution..`)
		time.Sleep(time.Second * 2)

		// Recupèration des données
		repository = dirInfo.Path()
		character = inputCharacter.Text

		// Action
		result, _ := SortFileInDirectory(repository, character, checked)

		// Update labels
		label2.SetText(`Résultat : Terminé..`)
		label3.SetText(`Gestion d'erreur : ` + result)
	})

	quitButton := widget.NewButton("Quitter", func() {
		myWindow.Close()
		os.Exit(0)
	})

	vBox1 := container.NewVBox(
		label,
		buttonFolder,
	)

	hBox2 := container.NewHBox(
		label1,
	)

	vBox3 := container.NewVBox(
		launchButton,
		label2,
		label3,
	)

	hBox4 := container.NewHBox(
		quitButton,
	)

	content := container.NewVBox(
		vBox1,
		hBox2,
		inputCharacter,
		check,
		vBox3,
		hBox4,
	)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}

func SortFileInDirectory(repository string, character string, isCheck bool) (string, bool) {
	// Command Prompt doesn't exist in %PATH%, so I need to use cmd.exe and in arguments use /c cd
	args := strings.Split("/c cd "+repository, " ")

	cmd := exec.Command("cmd.exe", args...)
	// Set _ to "out" to use the output in []bytes and to use the directory where has been executed the script
	_, errCmd := cmd.Output()
	if errCmd != nil {
		fmt.Print(errCmd)
	}

	// Vérify if user put data
	if repository == "" {
		return "Aucun dossier spécifié.", false
	}

	if character == "" {
		return "Aucun caractère spécifié.", false
	}

	// Read the directory to see files name which are into
	filesName, errRepository := os.ReadDir(repository)
	if errRepository != nil {
		return "La destination donné n'a pas été trouvé, elle n'existe peut-être pas.", false
	}

	// To find files which has charset given by the user
	var printFile string
	for _, file := range filesName {
		for i := 0; i < len(file.Name())-(len(character)-1); i++ {
			if len(character) < 2 {
				if file.Name()[i:i+len(character)] == character {
					printFile += file.Name() + "\n"
					break
				}
			} else {
				if file.Name()[i:i+len(character)] == character {
					printFile += file.Name() + "\n"
					break
				}
			}
		}
	}

	// If there is no files with the charset of the user
	if printFile == "" {
		return "Aucun fichier trouvé.", false
	}

	// Counter to know how much files there is
	var LenTab int
	for _, Char := range printFile {
		if Char == 10 {
			LenTab++
		}
	}

	// Append nothing to TabFile to init is length
	var TabFile []string
	for LenTab != 0 {
		TabFile = append(TabFile, "")
		LenTab--
	}

	// Append each name file into TabFile
	var count int
	for _, Char := range printFile {
		if Char == 10 {
			count++
		} else {
			TabFile[count] += string(Char)
		}
	}

	// Define repository
	var Path string
	if isCheck != true {
		Path = repository + "/" + "Vos_fichiers"
	} else {
		Path = os.Getenv("USERPROFILE") + "/" + "Desktop" + "/" + "Vos_fichiers"
	}

	// Create new repository named Vos_fichiers
	errCreateDir := os.Mkdir(Path, fs.ModePerm)
	if errCreateDir != nil {
		log.Println(errCreateDir)
		return "Impossible de créer le répertoire.", false
	}

	// Move files into new repository
	var errMoveFile error
	var NewPath string
	for i := 0; i < len(TabFile); i++ {
		OriginalPath := repository + "/" + TabFile[i]
		NewPath = Path + "/" + TabFile[i]

		errMoveFile = os.Rename(OriginalPath, NewPath)
		if errMoveFile != nil {
			return `Soucis dans le déplacement des fichiers.`, false
		}
	}

	return `Vos fichiers ont correctement été déplacés ! Dans le fichier "Vos_fichiers".`, true
}
