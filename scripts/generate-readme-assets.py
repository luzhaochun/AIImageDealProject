#!/usr/bin/env python3
"""生成 README 用 demo.gif 与 GitHub Social Preview 封面图。"""
from pathlib import Path

try:
    from PIL import Image, ImageDraw, ImageFont
except ImportError:
    raise SystemExit("请先安装: pip3 install Pillow")

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "docs" / "images" / "canvas-studio"
OUT.mkdir(parents=True, exist_ok=True)

SOFA_BEFORE = ROOT / "docs/images/canvas-studio/demo-sofa-before.png"
SOFA_AFTER = ROOT / "docs/images/canvas-studio/demo-sofa-after.png"
COWS_BEFORE = ROOT / "docs/images/canvas-studio/demo-cows-before.png"
COWS_AFTER = ROOT / "docs/images/canvas-studio/demo-cows-after.png"
UI_SHOT = ROOT / "docs/images/canvas-studio/ui-screenshot.png"


def fit(img: Image.Image, w: int, h: int) -> Image.Image:
    img = img.convert("RGB")
    img.thumbnail((w, h), Image.Resampling.LANCZOS)
    canvas = Image.new("RGB", (w, h), (30, 41, 59))
    ox = (w - img.width) // 2
    oy = (h - img.height) // 2
    canvas.paste(img, (ox, oy))
    return canvas


def make_demo_gif():
    frames = []
    fw, fh = 480, 270
    pairs = [
        (SOFA_BEFORE, SOFA_AFTER, "沙发场景 · 消除人物"),
        (COWS_BEFORE, COWS_AFTER, "牧场 · 消除背景牛群"),
    ]
    for before, after, _ in pairs:
        b = fit(Image.open(before), fw, fh)
        a = fit(Image.open(after), fw, fh)
        for _ in range(8):
            frames.append(b.copy())
        for i in range(6):
            frames.append(Image.blend(b, a, i / 5))
        for _ in range(8):
            frames.append(a.copy())
    out = OUT / "demo.gif"
    frames[0].save(
        out,
        save_all=True,
        append_images=frames[1:],
        duration=180,
        loop=0,
        optimize=True,
    )
    print("wrote", out, f"({len(frames)} frames)")


def make_social_preview():
    w, h = 1280, 640
    bg = (15, 23, 42)  # #0f172a
    canvas = Image.new("RGB", (w, h), bg)
    draw = ImageDraw.Draw(canvas)
    try:
        title_font = ImageFont.truetype("/System/Library/Fonts/Supplemental/Arial Bold.ttf", 42)
        sub_font = ImageFont.truetype("/System/Library/Fonts/Supplemental/Arial.ttf", 22)
        tag_font = ImageFont.truetype("/System/Library/Fonts/Supplemental/Arial.ttf", 18)
    except OSError:
        title_font = ImageFont.load_default()
        sub_font = title_font
        tag_font = title_font

    draw.text((48, 40), "ImageDeal", fill=(255, 255, 255), font=title_font)
    draw.text((48, 92), "GPT 图像消除 · Canvas 工作室", fill=(56, 189, 248), font=sub_font)
    draw.text((48, 128), "SaaS 图片协作 + OpenAI gpt-image-2 Inpainting", fill=(148, 163, 184), font=tag_font)

    bw, bh = 380, 280
    gap = 24
    y0 = 200
    x1 = 48
    x2 = x1 + bw + gap
    before = fit(Image.open(SOFA_BEFORE), bw, bh)
    after = fit(Image.open(SOFA_AFTER), bw, bh)
    canvas.paste(before, (x1, y0))
    canvas.paste(after, (x2, y0))
    draw.text((x1, y0 - 28), "Before", fill=(248, 113, 113), font=tag_font)
    draw.text((x2, y0 - 28), "After · GPT Image", fill=(74, 222, 128), font=tag_font)
    draw.rectangle([x1, y0, x1 + bw, y0 + bh], outline=(51, 65, 85), width=2)
    draw.rectangle([x2, y0, x2 + bw, y0 + bh], outline=(51, 65, 85), width=2)

    if UI_SHOT.exists():
        ui = Image.open(UI_SHOT).convert("RGB")
        ui.thumbnail((420, 260), Image.Resampling.LANCZOS)
        ux = w - ui.width - 40
        uy = h - ui.height - 36
        canvas.paste(ui, (ux, uy))
        draw.rectangle([ux, uy, ux + ui.width, uy + ui.height], outline=(51, 65, 85), width=2)

    out = ROOT / "docs/images/social-preview.png"
    canvas.save(out, "PNG", optimize=True)
    print("wrote", out)


if __name__ == "__main__":
    make_demo_gif()
    make_social_preview()
