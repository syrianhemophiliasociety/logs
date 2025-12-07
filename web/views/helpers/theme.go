package helpers

import "context"

func ThemeCtx(ctx context.Context) string {
	theme, ok := ctx.Value("theme-name").(string)
	if !ok {
		return "black"
	}

	return theme
}

func BackgroundImageCtx(ctx context.Context) string {
	switch ThemeCtx(ctx) {
	case "white":
		return "/assets/images/shs-bg-logo.webp"
	case "black":
		fallthrough
	default:
		return "/assets/images/shs-bg-logo-dark.webp"
	}
}
