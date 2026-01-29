#!/bin/sh
# Script para probar conectividad SMTP desde el contenedor Docker

echo "=== Probando conectividad SMTP ==="
echo "Servidor: server29.hostfactory.ch:465"
echo ""

# Método 1: Usar nc (netcat) si está disponible
echo "1. Probando con nc (netcat)..."
if command -v nc >/dev/null 2>&1; then
    nc -zv -w 5 server29.hostfactory.ch 465
    if [ $? -eq 0 ]; then
        echo "✓ Conexión exitosa con nc"
    else
        echo "✗ Conexión fallida con nc"
    fi
else
    echo "nc no está instalado"
fi

echo ""

# Método 2: Usar openssl para probar conexión TLS
echo "2. Probando con openssl s_client..."
if command -v openssl >/dev/null 2>&1; then
    echo | openssl s_client -connect server29.hostfactory.ch:465 -verify_return_error 2>&1 | head -5
    if [ $? -eq 0 ]; then
        echo "✓ Conexión TLS exitosa"
    else
        echo "✗ Conexión TLS fallida"
    fi
else
    echo "openssl no está instalado"
fi

echo ""

# Método 3: Usar timeout con /dev/tcp (bash)
echo "3. Probando con timeout y /dev/tcp..."
if timeout 5 bash -c 'cat < /dev/null > /dev/tcp/server29.hostfactory.ch/465' 2>/dev/null; then
    echo "✓ Puerto 465 está accesible"
else
    echo "✗ Puerto 465 NO está accesible"
    echo "  Esto puede indicar:"
    echo "  - Firewall bloqueando el puerto"
    echo "  - Red del contenedor sin acceso externo"
    echo "  - Servidor SMTP no disponible"
fi

echo ""
echo "=== Prueba completada ==="
