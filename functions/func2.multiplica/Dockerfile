# Imagen base
FROM python:3.8-slim

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el archivo Python al contenedor
COPY multiplica.py /app/multiplica.py

# Comando de entrada
ENTRYPOINT ["python3", "multiplica.py"]
