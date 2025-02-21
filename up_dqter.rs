use std::path::{Path, PathBuf};
use std::process::Command;
use std::thread;
use std::time::Duration;
use chrono::Local;
use git2::{Repository, Status};
use log::{error, info, warn, LevelFilter};
use log4rs::{
    append::file::FileAppender,
    config::{Appender, Config, Root},
    encode::pattern::PatternEncoder,
};
use serde::{Deserialize, Serialize};
use std::fs;
use std::error::Error;
use std::fmt;

// Configuración personalizada que se leerá desde config.toml
#[derive(Debug, Serialize, Deserialize)]
struct Config {
    repositories: Vec<String>,
    interval_hours: u64,
    commit_message: String,
}

// Error personalizado para el bot
#[derive(Debug)]
enum GitBotError {
    GitError(git2::Error),
    IoError(std::io::Error),
    InvalidPath(String),
    ConfigError(String),
}

impl fmt::Display for GitBotError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            GitBotError::GitError(e) => write!(f, "Error de Git: {}", e),
            GitBotError::IoError(e) => write!(f, "Error de IO: {}", e),
            GitBotError::InvalidPath(path) => write!(f, "Path inválido: {}", path),
            GitBotError::ConfigError(msg) => write!(f, "Error de configuración: {}", msg),
        }
    }
}

impl Error for GitBotError {}

impl From<git2::Error> for GitBotError {
    fn from(err: git2::Error) -> GitBotError {
        GitBotError::GitError(err)
    }
}

impl From<std::io::Error> for GitBotError {
    fn from(err: std::io::Error) -> GitBotError {
        GitBotError::IoError(err)
    }
}

// Inicializa el sistema de logging
fn setup_logging() -> Result<(), Box<dyn Error>> {
    let logfile = FileAppender::builder()
        .encoder(Box::new(PatternEncoder::new("{d} - {l} - {m}\n")))
        .build("git_bot.log")?;

    let config = Config::builder()
        .appender(Appender::builder().build("logfile", Box::new(logfile)))
        .build(Root::builder().appender("logfile").build(LevelFilter::Info))?;

    log4rs::init_config(config)?;
    Ok(())
}

// Lee la configuración desde el archivo config.toml
fn read_config() -> Result<Config, GitBotError> {
    let config_content = fs::read_to_string("config.toml")
        .map_err(|e| GitBotError::ConfigError(format!("No se pudo leer config.toml: {}", e)))?;
    
    toml::from_str(&config_content)
        .map_err(|e| GitBotError::ConfigError(format!("Error al parsear config.toml: {}", e)))
}

// Verifica si un repositorio tiene cambios
fn has_changes(repo: &Repository) -> Result<bool, git2::Error> {
    let statuses = repo.statuses(None)?;
    Ok(statuses.iter().any(|status| {
        let status = status.status();
        status.is_wt_modified() || status.is_wt_new()
    }))
}

// Realiza commit y push para un repositorio
fn commit_and_push(repo_path: &Path, commit_message: &str) -> Result<(), GitBotError> {
    if !repo_path.exists() {
        return Err(GitBotError::InvalidPath(repo_path.to_string_lossy().into()));
    }

    let repo = Repository::open(repo_path)?;
    
    if !has_changes(&repo)? {
        info!("No hay cambios en {}", repo_path.display());
        return Ok(());
    }

    // Añadir todos los cambios
    let mut index = repo.index()?;
    index.add_all(["*"].iter(), git2::IndexAddOption::DEFAULT, None)?;
    index.write()?;

    // Crear commit
    let tree_id = index.write_tree()?;
    let tree = repo.find_tree(tree_id)?;
    
    let signature = repo.signature()?;
    let parent_commit = repo.head()?.peel_to_commit()?;
    
    repo.commit(
        Some("HEAD"),
        &signature,
        &signature,
        commit_message,
        &tree,
        &[&parent_commit],
    )?;

    // Push usando git command
    let output = Command::new("git")
        .current_dir(repo_path)
        .args(&["push", "origin", "main"])
        .output()?;

    if !output.status.success() {
        let error_msg = String::from_utf8_lossy(&output.stderr);
        error!("Error al hacer push: {}", error_msg);
        return Err(GitBotError::GitError(git2::Error::from_str(&error_msg)));
    }

    info!("Commit y push exitosos para {}", repo_path.display());
    Ok(())
}

fn main() -> Result<(), Box<dyn Error>> {
    // Configurar logging
    setup_logging()?;
    info!("Iniciando Git Bot");

    // Leer configuración
    let config = read_config()?;
    info!("Configuración cargada: intervalo de {} horas", config.interval_hours);

    loop {
        for repo_path in &config.repositories {
            let path = PathBuf::from(repo_path);
            match commit_and_push(&path, &config.commit_message) {
                Ok(_) => info!("Procesamiento exitoso para {}", repo_path),
                Err(e) => error!("Error procesando {}: {}", repo_path, e),
            }
        }

        // Esperar hasta el próximo ciclo
        let wait_duration = Duration::from_secs(config.interval_hours * 3600);
        info!("Esperando {} horas hasta el próximo ciclo", config.interval_hours);
        thread::sleep(wait_duration);
    }
}