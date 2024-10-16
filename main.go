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
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio de trabajo:", err)
        os.Exit(1)
    }
    fmt.Println("Directorio de trabajo actual:", cwd)

    loadEnv() // Carga las variables de entorno desde el archivo .env

    repos := []string{"/ruta/al/repo1", "/ruta/al/repo2"} // Lista de repositorios a procesar

    // Itera sobre cada repositorio y lo procesa
    for _, repo := range repos {
        if err := processRepo(repo); err != nil {
            fmt.Printf("Error procesando el repo %s: %v\n", repo, err)
        }
    }
}

// Carga las variables de entorno desde un archivo .env
func loadEnv() {
    err := godotenv.Load("/ruta/completa/al/archivo/.env")
    if err != nil {
        fmt.Println("Error al cargar el archivo .env:", err)
        os.Exit(1)
    }
    fmt.Println("Archivo .env cargado correctamente")
}

// Procesa un repositorio específico
func processRepo(repoPath string) error {
    // Obtiene el usuario y correo de Git desde las variables de entorno
    gitUser := os.Getenv("GIT_USER")
    gitEmail := os.Getenv("GIT_EMAIL")

    fmt.Printf("Usuario de Git: %s\n", gitUser)
    fmt.Printf("Email de Git: %s\n", gitEmail)

    if gitUser == "" || gitEmail == "" {
        fmt.Println("Las variables GIT_USER o GIT_EMAIL no están establecidas")
        os.Exit(1)
    }

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
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error al configurar Git %s: %v\nSalida: %s\n", key, err, string(output))
        os.Exit(1)
    }
    fmt.Printf("Configurado Git %s a %s en %s\n", key, value, repoPath)
}

// Realiza un pull en el repositorio para obtener los últimos cambios
func pullRepo(repoPath string) {
    cmd := exec.Command("git", "-C", repoPath, "pull")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error al hacer pull en el repo %s: %v\nSalida: %s\n", repoPath, err, string(output))
        os.Exit(1)
    }
    fmt.Printf("Pull realizado en %s\n", repoPath)
}

// Actualiza el archivo README.md agregando una línea con la fecha y hora actual
func updateReadme(repoPath string) {
    currentTime := time.Now().Format("2006-01-02 15:04:05") // Formatea la fecha y hora actual
    appendText := fmt.Sprintf("\nHola mundo - %s", currentTime) // Texto a agregar
    // Comando para agregar el texto al final del README.md
    cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' >> %s/README.md", appendText, repoPath))
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error al actualizar README en el repo %s: %v\nSalida: %s\n", repoPath, err, string(output))
        os.Exit(1)
    }
    fmt.Printf("README.md actualizado en %s\n", repoPath)
}

// Hace commit de los cambios en el repositorio con un mensaje que incluye la fecha y hora
func commitChanges(repoPath string) {
    currentTime := time.Now().Format("2006-01-02 15:04:05") // Formatea la fecha y hora actual
    commitMsg := fmt.Sprintf("Add 'Hola mundo' con fecha %s", currentTime) // Mensaje de commit
    // Comando para hacer commit de todos los cambios con el mensaje especificado
    cmd := exec.Command("git", "-C", repoPath, "commit", "-am", commitMsg)
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error al hacer commit en el repo %s: %v\nSalida: %s\n", repoPath, err, string(output))
        os.Exit(1)
    }
    fmt.Printf("Commit realizado en %s\n", repoPath)
}

// Empuja los cambios al repositorio remoto
func pushChanges(repoPath string) {
    cmd := exec.Command("git", "-C", repoPath, "push")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Error al hacer push en el repo %s: %v\nSalida: %s\n", repoPath, err, string(output))
        os.Exit(1)
    }
    fmt.Printf("Push realizado en %s\n", repoPath)
}