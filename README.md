# SVG to PDF Converter (Go + Gin + pdfcpu)

A simple web service that converts an uploaded **SVG file** into a **PDF** using:

- `rsvg-convert` (for converting SVG to PNG)
- `pdfcpu` (for generating the final PDF)
- `Gin` (web framework for Go)
- HTML/JS frontend for uploading and downloading the PDF

---

## 📦 Features

- Upload `.svg` files via frontend or `curl`
- Converts to `.png` using `rsvg-convert`
- Embeds PNG into a PDF using `pdfcpu`
- Downloads the resulting `PDF` file

---

## 🛠️ Requirements

- Go 1.18+
- [rsvg-convert](https://wiki.gnome.org/Projects/LibRsvg) installed  
  (usually available via `librsvg2-bin` or similar)
