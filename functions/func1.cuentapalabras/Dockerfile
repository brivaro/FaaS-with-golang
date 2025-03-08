# Imagen base
FROM python:3.8-slim

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el archivo Python al contenedor
COPY cuenta_palabras.py /app/cuenta_palabras.py

# Comando de entrada
ENTRYPOINT ["python3", "cuenta_palabras.py"]
