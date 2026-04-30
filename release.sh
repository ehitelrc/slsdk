#!/bin/bash

# Validate that a version argument was provided
if [ -z "$1" ]; then
  echo "❌ Error: Debes especificar una versión para el tag."
  echo "💡 Uso: ./release.sh <version> [mensaje_commit]"
  echo "👉 Ejemplo: ./release.sh v0.1.0 \"Release inicial del SDK\""
  exit 1
fi

VERSION=$1
COMMIT_MSG=${2:-"Release $VERSION"}

echo "🚀 Iniciando proceso de release para $VERSION..."

# 1. Añadir cambios
echo "📦 Añadiendo cambios al stage..."
git add .

# 2. Hacer commit
echo "📝 Creando commit con el mensaje: '$COMMIT_MSG'"
git commit -m "$COMMIT_MSG" || echo "⚠️  No hay cambios nuevos para commitear. Continuando con el tag..."

# 3. Obtener rama actual y subir los cambios a GitHub
BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "☁️  Subiendo cambios de la rama '$BRANCH' a GitHub..."
git push origin $BRANCH

# 4. Crear el tag
echo "🏷️  Creando tag local $VERSION..."
# Eliminar el tag local si ya existe para poder recrearlo (opcional, útil si hubo un error)
git tag -d $VERSION 2>/dev/null
git tag $VERSION

# 5. Subir el tag a GitHub
echo "🚀 Subiendo tag $VERSION a GitHub..."
# Eliminar el tag remoto si existe (para forzar la actualización si lo estamos re-lanzando)
# git push origin :refs/tags/$VERSION 2>/dev/null
git push origin $VERSION

echo ""
echo "✅ ¡Release $VERSION completado exitosamente!"
echo ""
echo "Para usar esta versión en otro proyecto de Go, asegúrate de limpiar la caché de módulos y ejecutar:"
echo "--------------------------------------------------------"
echo "  GOPROXY=direct go get github.com/ehitelrc/slsdk@$VERSION"
echo "--------------------------------------------------------"
