# GoNamer

[![Go Report Card](https://goreportcard.com/badge/github.com/nouuu/gonamer)](https://goreportcard.com/report/github.com/nouuu/gonamer)
[![Go Reference](https://pkg.go.dev/badge/github.com/nouuu/gonamer.svg)](https://pkg.go.dev/github.com/nouuu/gonamer)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nouuu/gonamer)](https://golang.org/doc/devel/release.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Build Status](https://github.com/nouuu/gonamer/workflows/build/badge.svg)](https://github.com/nouuu/gonamer/actions?query=workflow%3Abuild)
[![Tests](https://github.com/nouuu/gonamer/workflows/tests/badge.svg)](https://github.com/nouuu/gonamer/actions?query=workflow%3Atests)
[![Lint](https://github.com/nouuu/gonamer/workflows/lint/badge.svg)](https://github.com/nouuu/gonamer/actions?query=workflow%3Alint)
[![Security](https://github.com/nouuu/gonamer/workflows/security/badge.svg)](https://github.com/nouuu/gonamer/actions?query=workflow%3Asecurity)

[![Release](https://img.shields.io/github/v/release/nouuu/gonamer)](https://github.com/nouuu/gonamer/releases)
[![Issues](https://img.shields.io/github/issues/nouuu/gonamer)](https://github.com/nouuu/gonamer/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/nouuu/gonamer)](https://github.com/nouuu/gonamer/pulls)
[![Contributors](https://img.shields.io/github/contributors/nouuu/gonamer)](https://github.com/nouuu/gonamer/graphs/contributors)
[![Lines of Code](https://tokei.rs/b1/github/nouuu/gonamer)](https://github.com/nouuu/gonamer)
[![Last Commit](https://img.shields.io/github/last-commit/nouuu/gonamer)](https://github.com/nouuu/gonamer/commits/main)

GoNamer is a powerful media file renaming tool that uses the TMDB API to automatically organize and rename your movie and TV show files based on accurate metadata.

## Features

- 🎬 Automatic movie and TV show recognition
- 🔍 Smart title detection with fuzzy matching
- 📺 Episode and season number detection
- 🎯 TMDB API integration for accurate metadata
- 🔄 Concurrent processing for better performance
- 💾 Integrated caching system
- 📝 Customizable naming patterns
- 🚀 Dry-run mode for safe testing
- 🌐 Multi-language support
- 🔒 Safe renaming with conflict prevention

## Installation

### Using Go

```bash
go install github.com/nouuu/gonamer/cmd@latest
```

### From Source

```bash
git clone https://github.com/nouuu/gonamer.git
cd gonamer
make install
```

### Package Managers

Coming soon:
- Homebrew
- APT
- RPM
- AUR

## Quick Start

1. Set your TMDB API key:
```bash
export TMDB_API_KEY=your_api_key
```

2. Run GoNamer:
```bash
gonamer /path/to/media
```

## Configuration

GoNamer can be configured using environment variables or a configuration file (coming soon):

```env
TMDB_API_KEY=your_key       # Required
MEDIA_PATH=./              # Path to scan
RECURSIVE=true            # Scan subdirectories
INCLUDE_NOT_FOUND=false  # Include unmatched files
DRY_RUN=true            # Test without renaming
MOVIE_PATTERN="{name} - {year}{extension}"
TVSHOW_PATTERN="{name} - {season}x{episode}{extension}"
TYPE=movie             # movie or tvshow
MAX_RESULTS=5         # Max suggestions per file
QUICK_MODE=false     # Skip confirmation
```

## Usage Examples

### Movies

```bash
# Scan movies with default pattern
gonamer /path/to/movies

# Custom pattern with dry run
DRY_RUN=true MOVIE_PATTERN="{name} ({year}){extension}" gonamer /movies
```

### TV Shows

```bash
# Scan TV shows
TYPE=tvshow gonamer /path/to/shows

# Custom episode pattern
TVSHOW_PATTERN="{name} - S{season}E{episode} - {episode_title}{extension}" gonamer /shows
```

## Development

### Prerequisites

- Go 1.22 or higher
- Make
- TMDB API key

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap 🗺️

### ✅ Completed Features
- [x] Basic CLI interface
- [x] TMDB API integration
- [x] Movie renaming support
- [x] TV show renaming support
- [x] Environment-based configuration
- [x] Cache system with Ristretto
- [x] Concurrent file processing
- [x] Custom naming patterns
- [x] Interactive mode
- [x] Dry-run support

### 🚀 Current Focus
- [x] Prepare for Github with CI implementation
    - Setup GitHub Actions workflows
    - Implement comprehensive testing
    - Add security checks
    - Configure automated releases
- [x] Migrate from .env to config.yml file
    - Design YAML configuration structure
    - Implement configs file loading
    - Add validation layer
- [x] Enhance CLI with Cobra framework
    - Migrate to Cobra commands
    - Add command documentation
    - Implement config path override
- [ ] Improve TV Show Processing
    - Implement folder-based processing
    - Add season folder structure support
    - Improve episode detection

### 🔮 Future Improvements
- [ ] Enhanced User Experience
    - Better suggestion selection UI
    - Progress bars for batch operations
    - Preview mode with detailed changes
    - Summary report after operations
     
- [ ] Web interface
  - A better UI with web interface

- [ ] Docker image ?

- [ ] Advanced File Management
    - Cache for previously processed files
    - Handling of subtitle files
    - Support for multi-episode files
    - Automatic creation of season folders

- [ ] Pattern System Enhancement
    - Visual pattern builder
    - More naming variables (quality, audio, etc.)
    - Per-folder pattern configuration
    - Pattern validation and testing

- [ ] API Integration
    - Support for additional APIs (IMDB, TVMaze)
    - Automatic API failover
    - Better metadata matching
    - Extended language support

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with love using Go
- Powered by [TMDB API](https://www.themoviedb.org/documentation/api)
- Inspired by the need for better media organization

---

Made with ❤️ by [nouuu](https://github.com/nouuu)
