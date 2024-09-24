package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

// Función principal que ejecuta el programa
func main() {
	loadEnv() // Carga las variables de entorno desde el archivo .env
	repos := []string{"/path/to/repo1", "/path/to/repo2"} // Lista de repositorios a procesar

	// Itera sobre cada repositorio y lo procesa
	for _, repo := range repos {
		if err := processRepo(repo); err != nil {
			fmt.Printf("Error procesando el repo %s: %v\n", repo, err)
		}
	}
}

// Carga las variables de entorno desde un archivo .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error al cargar el archivo .env")
		os.Exit(1) // Termina el programa si no se puede cargar el archivo .env
	}
}

// Procesa un repositorio específico
func processRepo(repoPath string) error {
	// Obtiene el usuario y correo de Git desde las variables de entorno
	gitUser := os.Getenv("GIT_USER")
	gitEmail := os.Getenv("GIT_EMAIL")

	// Configura el usuario y correo en el repositorio
	setGitConfig(repoPath, "user.name", gitUser)
	setGitConfig(repoPath, "user.email", gitEmail)

	// Realiza un pull para actualizar el repositorio
	pullRepo(repoPath)
	// Actualiza el archivo README.md
	updateReadme(repoPath)
	// Hace commit de los cambios
	commitChanges(repoPath)
	// Empuja los cambios al repositorio remoto
	pushChanges(repoPath)

	return nil
}

// Configura una clave de Git en el repositorio especificado
func setGitConfig(repoPath, key, value string) {
	cmd := exec.Command("git", "-C", repoPath, "config", key, value)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error al configurar Git %s: %v\n", key, err)
		os.Exit(1) // Termina el programa si falla la configuración
	}
}

// Realiza un pull en el repositorio para obtener los últimos cambios
func pullRepo(repoPath string) {
	cmd := exec.Command("git", "-C", repoPath, "pull")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error al hacer pull en el repo %s: %v\n", repoPath, err)
		os.Exit(1) // Termina el programa si falla el pull
	}
}

// Actualiza el archivo README.md agregando una línea con la fecha y hora actual
func updateReadme(repoPath string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05") // Formatea la fecha y hora actual
	appendText := fmt.Sprintf("\nHola mundo - %s", currentTime) // Texto a agregar
	// Comando para agregar el texto al final del README.md
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' >> %s/README.md", appendText, repoPath))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error al actualizar README en el repo %s: %v\n", repoPath, err)
		os.Exit(1) // Termina el programa si falla la actualización
	}
}

// Hace commit de los cambios en el repositorio con un mensaje que incluye la fecha y hora
func commitChanges(repoPath string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05") // Formatea la fecha y hora actual
	commitMsg := fmt.Sprintf("Add 'Hola mundo' con fecha %s", currentTime) // Mensaje de commit
	// Comando para hacer commit de todos los cambios con el mensaje especificado
	cmd := exec.Command("git", "-C", repoPath, "commit", "-am", commitMsg)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error al hacer commit en el repo %s: %v\n", repoPath, err)
		os.Exit(1) // Termina el programa si falla el commit
	}
}

// Empuja los cambios al repositorio remoto
func pushChanges(repoPath string) {
	cmd := exec.Command("git", "-C", repoPath, "push")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error al hacer push en el repo %s: %v\n", repoPath, err)
		os.Exit(1) // Termina el programa si falla el push
	}
}