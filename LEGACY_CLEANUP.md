# Legacy Files - Safe to Delete

The following files/folders are from the old project structure and have been migrated to the new organized structure. They can be safely deleted:

## Legacy Folders (duplicate content now in `backend/`)
- `cmd/` → Migrated to `backend/cmd/`
- `db/` → Migrated to `backend/db/`
- `internal/` → Migrated to `backend/internal/`
- `models/` → Migrated to `backend/models/`

## Legacy Files
- `go.mod` → Migrated to `backend/go.mod`
- `go.sum` → Migrated to `backend/go.sum`
- `internal/index.html` → Migrated to `frontend/html/index.html`

## Files to Keep

### Root Level
- `README.md` - Project documentation (updated)
- `commands.txt` - Command reference
- `text.txt` - Notes

### Backend Folder
- `backend/cmd/main.go` - **Server entry point** (ACTIVE)
- `backend/db/db.go` - **Database connection** (ACTIVE)
- `backend/internal/` - **API handlers** (ACTIVE)
- `backend/models/` - **Data models** (ACTIVE)
- `backend/go.mod` - **Dependencies** (ACTIVE)
- `backend/go.sum` - **Dependency checksums** (ACTIVE)

### Frontend Folder
- `frontend/html/index.html` - **Main page** (ACTIVE)
- `frontend/css/style.css` - **Styling** (ACTIVE)
- `frontend/js/script.js` - **Client logic** (ACTIVE)

## Cleanup Commands

To remove the legacy files:

```bash
cd /home/dev-diego/Desktop/VSCode/MongoDB-Project

# Remove old folders
rm -rf cmd/ db/ internal/ models/
rm -f go.mod go.sum
```

After cleanup, all active code will be in:
- `backend/` - All Go backend code
- `frontend/` - All frontend files (HTML, CSS, JS)
