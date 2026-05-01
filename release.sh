#!/bin/bash

# 1. Extraer versión actual del README.md
CURRENT_VERSION=$(grep -oE 'github\.com\/ehitelrc\/slsdk@v[0-9]+\.[0-9]+\.[0-9]+' README.md | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1)

if [ -z "$CURRENT_VERSION" ]; then
  echo "❌ Error: No se pudo encontrar una versión actual en el README.md."
  exit 1
fi

echo "📌 Versión actual en README.md: $CURRENT_VERSION"

# 2. Determinar la nueva versión
if [ -z "$1" ]; then
  # Si no se envía parámetro, autoincrementar el parche (PATCH)
  BASE_VER=${CURRENT_VERSION#v}
  IFS='.' read -ra PARTS <<< "$BASE_VER"
  MAJOR=${PARTS[0]}
  MINOR=${PARTS[1]}
  PATCH=${PARTS[2]}
  NEXT_PATCH=$((PATCH + 1))
  NEXT_VERSION="v${MAJOR}.${MINOR}.${NEXT_PATCH}"
  echo "💡 No se especificó versión explícita. Infiriendo siguiente parche: $NEXT_VERSION"
else
  NEXT_VERSION=$1
fi

COMMIT_MSG=${2:-"Release $NEXT_VERSION"}

echo "🚀 Iniciando proceso de release para $NEXT_VERSION..."

# 3. Actualizar README.md antes del commit
echo "📄 Actualizando referencias de versión en README.md..."
sed -i '' -E "s/github.com\/ehitelrc\/slsdk@[vV][0-9]+\.[0-9]+\.[0-9]+/github.com\/ehitelrc\/slsdk@$NEXT_VERSION/g" README.md

# 4. Añadir cambios
echo "📦 Añadiendo cambios al stage..."
git add .

# 5. Hacer commit
echo "📝 Creando commit con el mensaje: '$COMMIT_MSG'"
git commit -m "$COMMIT_MSG" || echo "⚠️  No hay cambios nuevos para commitear. Continuando con el tag..."

# 6. Obtener rama actual y subir los cambios a GitHub
BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "☁️  Subiendo cambios de la rama '$BRANCH' a GitHub..."
git push origin $BRANCH

# 7. Crear el tag
echo "🏷️  Creando tag local $NEXT_VERSION..."
git tag -d $NEXT_VERSION 2>/dev/null
git tag $NEXT_VERSION

# 8. Subir el tag a GitHub
echo "🚀 Subiendo tag $NEXT_VERSION a GitHub..."
git push origin $NEXT_VERSION

echo ""
echo "✅ ¡Release $NEXT_VERSION completado exitosamente!"
echo ""
echo "Para usar esta versión en otro proyecto de Go, asegúrate de limpiar la caché de módulos y ejecutar:"
echo "--------------------------------------------------------"
echo "  GOPROXY=direct go get github.com/ehitelrc/slsdk@$NEXT_VERSION"
echo "--------------------------------------------------------"
