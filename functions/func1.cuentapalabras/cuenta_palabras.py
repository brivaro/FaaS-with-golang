import json, sys

def cuenta_palabras(json_input):
    """
    Procesa un JSON, extrae el valor de 'frase' y cuenta el número total de palabras.

    Args:
        json_input (str): Cadena JSON con la clave 'frase'.

    Returns:
        int: Número total de palabras en la frase.
    """
    datos = json.loads(json_input)
    
    texto = datos.get("frase", "")
    palabras = texto.split()
    
    return (len(palabras))

if __name__ == "__main__":
    # Captura el primer argumento de la línea de comandos como el texto de entrada
    json_input = sys.argv[1]
    # Ejecuta la función y devuelve el resultado
    print(cuenta_palabras(json_input))  # Imprime el resultado en stdout

