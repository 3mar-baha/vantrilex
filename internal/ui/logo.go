package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// loadLogo renders the compact header crest via chafa at 28x14 (strict 2:1
// terminal-cell scaling) so the full emblem including both outer horns is
// always visible inside menu screens.
func loadLogo() string {
	return loadLogoSized(28, 14)
}

// loadLogoSized renders logo_crest.png via chafa at a controlled WxH size.
//
// 2:1 Terminal Cell Rule: character cells are ~twice as tall as wide, so the
// width argument must equal exactly 2x the height argument for a true 1:1
// undistorted render (hero 36x18, header 28x14).
//
// NOTE: chafa has no "--bg none" mode (it rejects the value); transparent
// pixels of the PNG fall through to the terminal background by default,
// which is exactly the desired effect, so no --bg flag is passed.
func loadLogoSized(w, h int) string {
	candidates := []string{
		"logo_crest.png",
		"logo.jpeg",
		"logo.jpg",
		filepath.Join("..", "logo_crest.png"),
		filepath.Join("..", "logo.jpeg"),
		filepath.Join("..", "..", "logo_crest.png"),
		filepath.Join("..", "..", "logo.jpeg"),
	}
	// Also try exe directory and cwd.
	if exe, err := os.Executable(); err == nil {
		d := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(d, "logo_crest.png"),
			filepath.Join(d, "logo.jpeg"),
			filepath.Join(d, "logo.jpg"),
		)
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "logo_crest.png"),
			filepath.Join(cwd, "logo.jpeg"),
		)
		// Walk up two levels.
		p := cwd
		for i := 0; i < 2; i++ {
			p = filepath.Dir(p)
			candidates = append(candidates,
				filepath.Join(p, "logo_crest.png"),
				filepath.Join(p, "logo.jpeg"),
			)
		}
	}
	var found string
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			found = c
			break
		}
	}
	if found == "" {
		return fallbackLogo()
	}
	if _, err := exec.LookPath("chafa"); err != nil {
		return fallbackLogo()
	}
	if w < 20 {
		w = 20
	}
	if w > 72 {
		w = 72
	}
	if h < 4 {
		h = 4
	}
	if h > 24 {
		h = 24
	}
	cmd := exec.Command("chafa", "--format", "symbols", "--symbols", "block+border+space-wide-inverted",
		"--size", itoa(w)+"x"+itoa(h), "--stretch", found)
	out, err := cmd.CombinedOutput()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		return fallbackLogo()
	}
	art := strings.TrimRight(string(out), "\n")
	return haloStyle.Render(art)
}

func fallbackLogo() string {
	lines := []string{
		"██╗   ██╗ █████╗ ███╗   ██╗████████╗██████╗ ██╗██╗     ███████╗██╗  ██╗",
		"██║   ██║██╔══██╗████╗  ██║╚══██╔══╝██╔══██╗██║██║     ██╔════╝╚██╗██╔╝",
		"██║   ██║███████║██╔██╗ ██║   ██║   ██████╔╝██║██║     █████╗   ╚███╔╝ ",
		"╚██╗ ██╔╝██╔══██║██║╚██╗██║   ██║   ██╔══██╗██║██║     ██╔══╝   ██╔██╗ ",
		" ╚████╔╝ ██║  ██║██║ ╚████║   ██║   ██║  ██║██║███████╗███████╗██╔╝ ██╗",
		"  ╚═══╝  ╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚═╝  ╚═╝",
	}
	return titleStyle.Render(strings.Join(lines, "\n"))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [16]byte
	pos := len(b)
	for n > 0 {
		pos--
		b[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}
