import os
import re

import unicodedata

def remove_emojis(text):
    # Remove emojis and variation selectors
    text = "".join(c for c in text if unicodedata.category(c) not in ('So', 'Cn'))
    # Remove variation selectors
    text = re.sub(r'[\ufe00-\ufe0f]', '', text)
    return text

def remove_comments(code):
    # Removing multi-line comments
    code = re.sub(re.compile(r"/\*.*?\*/", re.DOTALL), "", code)
    
    # Removing single-line comments
    # We need to be careful not to remove // inside strings.
    # A simple way is to use a regex that matches strings first.
    def replacer(match):
        s = match.group(0)
        if s.startswith('/'):
            return "" # it's a comment
        else:
            return s # it's a string or something else
            
    pattern = re.compile(
        r'//.*?$|/\*.*?\*/|\'(?:\\.|[^\\\'])\'|"(?:\\.|[^\\"])*"',
        re.DOTALL | re.MULTILINE
    )
    return re.sub(pattern, replacer, code)

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Remove comments first
    content = remove_comments(content)
    # Remove emojis
    content = remove_emojis(content)
    
    # Clean up multiple empty lines
    content = re.sub(r'\n\s*\n\s*\n', '\n\n', content)
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content.strip() + '\n')

def main():
    root_dir = "/home/ramji/Desktop/reset/final-tick/Backend(Go)"
    for root, dirs, files in os.walk(root_dir):
        for file in files:
            if file.endswith(".go"):
                fullpath = os.path.join(root, file)
                print(f"Processing: {fullpath}")
                process_file(fullpath)

if __name__ == "__main__":
    main()
