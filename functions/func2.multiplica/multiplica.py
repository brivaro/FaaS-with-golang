import json, sys

def multiplica(json_input):
    """
    Extrae los números 'a' y 'b' de un JSON y retorna su producto.

    Args:
        json_input (str): Cadena JSON con las claves 'a' y 'b'.

    Returns:
        int/float: Producto de 'a' y 'b'.
    """
    datos = json.loads(json_input)
    
    a = datos.get("a", 0)
    b = datos.get("b", 0)
    
    return (a * b)

if __name__ == "__main__":
    # Captura el primer argumento de la línea de comandos como el texto de entrada
    json_input = sys.argv[1]
    # Ejecuta la función y devuelve el resultado
    print(multiplica(json_input))  # Imprime el resultado en stdout
