Plan de Implementación: Comando `anx-agent analyze`

**1. Objetivo del Comando `analyze`**

El comando `anx-agent analyze` tiene como propósito principal analizar el código fuente de un proyecto de software dado una ruta específica. Su objetivo es:
*   Extraer información clave sobre la estructura del proyecto, tecnologías utilizadas, propósito general y dependencias.
*   Generar una visión general inteligible del proyecto, incluyendo un resumen de alto nivel, puntos clave y posibles áreas de mejora.
*   Proporcionar la capacidad de generar reportes en diferentes formatos (e.g., Markdown, JSON) para facilitar su consumo y posterior procesamiento.
*   Manejar proyectos de gran tamaño dividiendo el código en "chunks" manejables para los modelos de IA, respetando los límites de tokens.

**2. Alcance (MVP)**

Para la primera iteración (MVP) del comando `analyze`, se incluirán las siguientes funcionalidades:
*   **Análisis de Ruta:** Soporte para analizar un directorio completo o un archivo individual especificado por la flag `--path`.
*   **Filtrado de Archivos:** Capacidad de incluir/excluir archivos basándose en sus extensiones (flag `--extensions`).
*   **Ignorar Rutas:** Soporte para ignorar directorios o archivos específicos (e.g., `node_modules`, `.git`, `bin`) usando una flag `--ignore-paths`.
*   **Procesamiento de Contenido:** Lectura y preparación del contenido de archivos de texto/código.
*   **Chunking Básico:** Estrategia inicial de división de archivos grandes en partes más pequeñas para adecuarse a los límites de tokens de la IA.
*   **Integración AI:** Utilización del modelo de IA configurado (Gemini por defecto) para realizar el análisis.
*   **Generación de Reportes:** Salida del análisis en formato Markdown y JSON, especificado por la flag `--output`.

**3. Diseño de Alto Nivel**

El comando `analyze` se integrará en la arquitectura existente de ANX Agent, aprovechando los módulos ya definidos:
*   **`cmd/agentcli`**: Punto de entrada principal para la definición del comando `analyze` y la gestión de sus flags.
*   **`internal/cli/analyze`**: Contendrá la lógica específica para la ejecución del comando, orquestando las llamadas a otros módulos.
*   **`internal/agent/project`**: Módulo encargado del descubrimiento y filtrado de archivos dentro de la estructura del proyecto.
*   **`internal/agent/chunker`**: Nuevo módulo responsable de dividir el contenido de archivos grandes en fragmentos (chunks) aptos para la IA.
*   **`internal/ai`**: Abstracción para la interacción con los proveedores de IA (inicialmente Gemini). Se utilizará para enviar los chunks y recibir análisis parciales.
*   **`internal/agent/analyzer`**: Nuevo módulo que contendrá la lógica principal para orquestar el análisis, desde la preparación de prompts hasta la consolidación de las respuestas de la IA.
*   **`internal/reporting`**: Módulo ya existente que se extenderá para formatear el resultado final del análisis en los formatos solicitados.
*   **`internal/config`**: Usado para cargar y acceder a la configuración global (API keys, log levels, etc.).

```
┌─────────────────┐       ┌──────────────────┐
│   anx-agent CLI │       │ internal/config  │
│  (cmd/agentcli) ├───────► (Config Load)    │
└─────────────────┘       └──────────────────┘
        │
        ▼ Define Command
┌─────────────────────┐
│ internal/cli/analyze│
│ (Command Handler)   │
└─────────────────────┘
        │
        ▼ Orchestrates
┌───────────────────────────┐
│ internal/agent/project    │
│ (File Discovery & Filter) ├───────► List of `File` objects
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐
│ internal/agent/chunker    │
│ (Code Chunking Strategy)  ├───────► List of `Chunk` objects
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐       ┌──────────────────┐
│ internal/agent/analyzer   │       │ internal/ai/gemini│
│ (AI Prompting & Analysis) ├───────► (AI API Calls)   │
└───────────────────────────┘       └──────────────────┘
        │
        ▼ `ProjectAnalysis` struct
┌───────────────────────────┐
│ internal/reporting        │
│ (Report Generation)       │
└───────────────────────────┘
        │
        ▼
┌──────────────────┐
│ Report Output    │
└──────────────────┘
```

**4. Pasos Detallados de Implementación**

**4.1. Definición del Comando CLI (`cmd/agentcli/cmd/analyze.go`)**
*   **Creación:** Añadir un nuevo archivo `analyze.go` dentro de `cmd/agentcli/cmd/`.
*   **`init()` function:** Registrar `analyzeCmd` como subcomando de `rootCmd`.
*   **`analyzeCmd` Definition:**
    *   `Use: "analyze"`
    *   `Short: "Analyze a project or specific files with AI"`
    *   `Long: "The analyze command scans a given directory or file path, processes its content using AI, and provides a comprehensive overview or report."`
    *   `RunE`: Se asignará una función `RunAnalyzeCmd` que contendrá la lógica principal.
*   **Flags:**
    *   `--path` (string, `viper.BindPFlag`): La ruta al proyecto o archivo a analizar. Hacerlo requerido.
    *   `--output` (string, default "markdown"): Formato del reporte (e.g., "markdown", "json").
    *   `--extensions` (string slice): Lista de extensiones de archivo a incluir (e.g., "go,md,txt").
    *   `--ignore-paths` (string slice): Lista de patrones de rutas a ignorar (e.g., "node_modules", ".git", "vendor", "test").
    *   `--max-file-size` (int, default X): Tamaño máximo de archivo en KB para incluir en el análisis. Archivos más grandes serán ignorados o procesados de manera diferente.
    *   `--prompt-template` (string): Ruta a un archivo de plantilla de prompt personalizado. (Considerar para futuras iteraciones si el tiempo lo permite).

**4.2. Lógica del Comando (`internal/cli/analyze/handler.go`)**
*   **Función `RunAnalyzeCmd(cmd *cobra.Command, args []string) error`:**
    *   Obtener los valores de las flags (`path`, `output`, `extensions`, `ignorePaths`, `maxFileSize`).
    *   Inicializar la configuración (`config.LoadConfig()`).
    *   Inicializar el cliente de AI (`ai.NewAIClient()`, usando el proveedor configurado, ej. Gemini).
    *   **Paso 1: Descubrimiento de Archivos:** Llamar a `project.DiscoverFiles(path, extensions, ignorePaths, maxFileSize)`.
        *   Manejar el caso en que `--path` sea un archivo único.
    *   **Paso 2: Análisis del Proyecto:** Instanciar un `analyzer.ProjectAnalyzer` y llamar a `analyzer.AnalyzeProject(discoveredFiles, aiClient)`.
        *   Este método orquestará el chunking y las llamadas a la IA.
    *   **Paso 3: Generación de Reporte:** Llamar a `reporting.GenerateReport(analysisResult, outputFormat)`.
    *   **Paso 4: Salida:** Imprimir el reporte generado a `stdout` o guardarlo en un archivo si se especifica una flag de salida (`--output-file`).

**4.3. Descubrimiento y Filtrado de Archivos (`internal/agent/project/finder.go`)**
*   **Estructura `File`:**
    ```go
    type File struct {
        Path    string // Ruta completa del archivo
        Content []byte // Contenido del archivo
        Ext     string // Extensión del archivo
    }
    ```
*   **Función `DiscoverFiles(rootPath string, extensions []string, ignorePaths []string, maxFileSizeKB int) ([]File, error)`:**
    *   Utilizar `filepath.Walk` para recorrer `rootPath` recursivamente.
    *   Para cada archivo encontrado:
        *   **Filtrado por `ignorePaths`:** Implementar lógica para saltar directorios o archivos que coincidan con los patrones de `ignorePaths`. Usar `path/filepath.Match` o similar para patrones.
        *   **Filtrado por `extensions`:** Si se especifican extensiones, verificar que el archivo tenga una de ellas.
        *   **Filtrado por `maxFileSizeKB`:** Si el tamaño del archivo excede `maxFileSizeKB`, ignorarlo o registrar una advertencia.
        *   Leer el contenido del archivo y crear una instancia `File`.
    *   Retornar la lista de `File`s.

**4.4. Estrategia de Chunking de Código (`internal/agent/chunker/code_chunker.go`)**
*   **Problema:** Los modelos de IA tienen límites de tokens (context window). Un archivo grande debe ser dividido.
*   **Estrategia MVP (Line-based Chunking):**
    *   **`ChunkContent(content []byte, maxTokensPerChunk int, tokenizer func(string) int) ([]Chunk, error)`:**
        *   Recibir el contenido de un archivo y un límite de tokens por chunk.
        *   Utilizar un `tokenizer` (e.g., un wrapper para el tokenizer de Gemini) para estimar los tokens.
        *   Dividir el contenido por líneas o bloques.
        *   Acumular líneas hasta que el número de tokens se acerque al `maxTokensPerChunk`.
        *   Cada segmento acumulado se convierte en un `Chunk`.
        *   **`Chunk` Structure:**
            ```go
            type Chunk struct {
                Content    string
                Metadata   map[string]string // e.g., "filename", "line_start", "line_end"
                TokenCount int
            }
            ```
*   **Consideraciones para el Chunking:**
    *   Intentar no romper bloques de código lógicos (e.g., funciones, clases). Para un MVP, una división por líneas es suficiente, pero se puede mejorar buscando patrones.
    *   Añadir metadatos al chunk para que la IA sepa de qué archivo y qué parte viene.

**4.5. Interacción con la IA (`internal/ai/gemini/client.go` y `internal/agent/analyzer/analyzer.go`)**
*   **`internal/ai/gemini/client.go`:**
    *   Asegurar la implementación del método `GenerateContent(prompt string) (string, error)` o similar.
    *   Manejo de reintentos y timeouts según la configuración de Viper.
    *   (Opcional) Implementar un método para estimar tokens en un string (`CountTokens(text string) (int, error)`).
*   **`internal/agent/analyzer/analyzer.go`:**
    *   **Estructura `ProjectAnalysis`:** Representará el resultado final del análisis.
        ```go
        type ProjectAnalysis struct {
            Overview        string         `json:"overview"` // Resumen general del proyecto
            KeyTechnologies []string       `json:"key_technologies"` // Tecnologías detectadas
            Structure       string         `json:"structure"` // Descripción de la estructura del directorio
            FileSummaries   []FileSummary  `json:"file_summaries"` // Resúmenes individuales de archivos
            Suggestions     []string       `json:"suggestions"` // Sugerencias de mejora
            PotentialIssues []string       `json:"potential_issues"` // Problemas detectados
        }

        type FileSummary struct {
            FilePath      string `json:"file_path"`
            Summary       string `json:"summary"` // Resumen del archivo
            KeyFunctions  []string `json:"key_functions,omitempty"`
            Dependencies  []string `json:"dependencies,omitempty"`
        }
        ```
    *   **`AnalyzeProject(files []project.File, aiClient ai.AIClient) (*ProjectAnalysis, error)`:**
        *   **Fase 1: Análisis por Archivo/Chunk:**
            *   Iterar sobre cada `project.File`.
            *   Para cada `File`:
                *   Generar `Chunks` usando `chunker.ChunkContent`.
                *   Para cada `Chunk`:
                    *   Construir un prompt específico para el chunk (e.g., "Analyze this part of file X. What is its purpose?"). Incluir el `chunk.Content` y `chunk.Metadata`.
                    *   Enviar el prompt al `aiClient.GenerateContent()`.
                    *   Recopilar la respuesta del AI.
                *   Si un archivo es pequeño y no necesita chunking, enviar su contenido completo con un prompt de "análisis de archivo".
        *   **Fase 2: Consolidación del Análisis:**
            *   Una vez que se tienen los análisis individuales de todos los archivos/chunks, construir un prompt de consolidación.
            *   **Prompt Consolidación Ejemplo:** "Based on the following summaries of files from a project, provide a high-level overview, identify key technologies, describe the project structure, and suggest improvements. Summaries:\n[Concatenated `FileSummary.Summary` and other relevant details from individual analyses]"
            *   Enviar este prompt al `aiClient.GenerateContent()` para obtener el `ProjectAnalysis.Overview`, `KeyTechnologies`, `Suggestions`, etc.
            *   Parsear las respuestas del AI para poblar la estructura `ProjectAnalysis`.

**4.6. Generación de Reportes (`internal/reporting/reporter.go`)**
*   **`GenerateReport(analysis *analyzer.ProjectAnalysis, format string) ([]byte, error)`:**
    *   **Markdown Formatter:** Convertir la estructura `ProjectAnalysis` a un string Markdown bien formateado, usando encabezados, listas, bloques de código, etc.
        *   Ejemplo de estructura Markdown:
            ```markdown
            # ANX Agent Project Analysis Report

            ## Project Overview
            [analysis.Overview]

            ## Key Technologies
            - [tech1]
            - [tech2]

            ## Project Structure
            [analysis.Structure]

            ## File Summaries
            ### `path/to/file1.go`
            [file1.Summary]
            Key Functions: [func1, func2]

            ### `path/to/file2.md`
            [file2.Summary]

            ## Suggestions & Potential Issues
            - [suggestion1]
            - [issue1]
            ```
    *   **JSON Formatter:** Serializar la estructura `ProjectAnalysis` a JSON usando `json.MarshalIndent`.

**4.7. Manejo de Errores y Logging**
*   Implementar manejo robusto de errores en todas las fases (I/O de archivos, errores de red/API de AI, errores de parsing).
*   Utilizar el sistema de logging de Go (o una librería como `logrus` si se integra) con el `log_level` configurado en Viper para proporcionar retroalimentación útil durante la ejecución (e.g., archivos ignorados, progreso, errores de AI).

**5. Estructuras de Datos / Interfaces Clave**

*   `internal/agent/project.File`: Representa un archivo descubierto.
*   `internal/agent/chunker.Chunk`: Representa un fragmento de contenido para el análisis AI.
*   `internal/ai.AIClient` interface: Abstracción para los clientes de IA.
*   `internal/agent/analyzer.ProjectAnalysis`: Resultado consolidado del análisis del proyecto.
*   `internal/agent/analyzer.FileSummary`: Resumen de un archivo individual.

**6. Estrategia de Pruebas**

*   **Pruebas Unitarias:**
    *   `internal/agent/project/finder_test.go`: Pruebas para el descubrimiento y filtrado de archivos (`DiscoverFiles`). Usar `os.TempDir` para crear estructuras de directorios de prueba.
    *   `internal/agent/chunker/code_chunker_test.go`: Pruebas para la lógica de chunking, asegurando que los chunks no excedan el tamaño y contengan metadatos correctos.
    *   `internal/ai/gemini/client_test.go`: Pruebas para la interacción con la API de Gemini (usando mocks para la API real).
    *   `internal/reporting/reporter_test.go`: Pruebas para la generación de reportes en diferentes formatos a partir de una estructura `ProjectAnalysis` mock.
*   **Pruebas de Integración:**
    *   `internal/agent/analyzer/analyzer_test.go`: Probar el flujo completo de análisis (chunking, llamadas simuladas a IA, consolidación).
*   **Pruebas E2E (End-to-End):**
    *   `cmd/agentcli/cli_test.go`: Ejecutar el comando `anx-agent analyze` con diferentes flags y `testdata` (directorio `testdata/`) y verificar la salida a `stdout` o archivos de salida.

**7. Mejoras Futuras**

*   **Chunking Inteligente:** En lugar de solo por líneas/tokens, implementar chunking basado en la estructura de AST (Abstract Syntax Tree) para lenguajes como Go, Python, etc., para enviar bloques de código más coherentes (funciones, clases, métodos) a la IA.
*   **Cache de AI:** Cachear respuestas de la API de AI para archivos que no han cambiado, reduciendo costos y tiempo de ejecución en análisis repetidos.
*   **Análisis Profundo:** Extender el análisis para incluir:
    *   Análisis de dependencias de módulos/paquetes.
    *   Posibles vulnerabilidades de seguridad (integración con herramientas SAST básicas).
    *   Análisis de complejidad de código (e.g., Ciclomática).
*   **Perfiles de Análisis:** Permitir al usuario definir "perfiles" de análisis que enfaticen diferentes aspectos (e.g., "seguridad", "rendimiento", "documentación").
*   **Plantillas de Prompts Configurables:** Permitir al usuario proporcionar sus propias plantillas de prompts para el análisis y consolidación.
*   **Soporte de `anxignore`:** Implementar un archivo de configuración `.anxignore` (similar a `.gitignore`) para especificar archivos y directorios a ignorar.
*   **Output Interactivo:** Proporcionar una salida más interactiva o visual para el reporte en el terminal.