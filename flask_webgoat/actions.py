import pickle
import base64
from pathlib import Path
import subprocess

from flask import Blueprint, request, jsonify, session

bp = Blueprint("actions", __name__)


@bp.route("/message", methods=["POST"])
def log_entry():
    user_info = session.get("user_info", None)
    if user_info is None:
        return jsonify({"error": "no user_info found in session"})
    access_level = user_info[2]
    if access_level > 2:
        return jsonify({"error": "access level < 2 is required for this action"})
    filename_param = request.form.get("filename")
    if filename_param is None:
        return jsonify({"error": "filename parameter is required"})
    text_param = request.form.get("text")
    if text_param is None:
        return jsonify({"error": "text parameter is required"})

    user_id = user_info[0]
    user_dir = "data/" + str(user_id)
    user_dir_path = Path(user_dir)
    if not user_dir_path.exists():
        user_dir_path.mkdir()

    filename = filename_param + ".txt"
    path = Path(user_dir + "/" + filename)
    with path.open("w", encoding="utf-8") as open_file:
        # vulnerability: Directory Traversal
        open_file.write(text_param)
    return jsonify({"success": True})


@bp.route("/grep_processes")
def grep_processes():
    name = request.args.get("name")
    
    # FIX: Validate that name parameter is provided and is a string
    if not name or not isinstance(name, str):
        return jsonify({"error": "invalid name parameter"})
    
    # FIX: Limit length to prevent DoS attacks
    if len(name) > 256:
        return jsonify({"error": "name parameter too long"})
    
    # FIX: Execute ps command directly without shell to eliminate command injection vulnerability
    try:
        res = subprocess.run(
            ["ps", "aux"],
            capture_output=True,
            timeout=5,
            text=True
        )
    except subprocess.TimeoutExpired:
        return jsonify({"error": "command timeout"})
    
    # FIX: Check if command executed successfully
    if res.returncode != 0:
        return jsonify({"error": "ps command failed"})
    
    # FIX: Parse output in Python instead of using shell pipes (grep/awk) to prevent command injection
    lines = res.stdout.strip().split('\n')
    names = []
    
    for line in lines[1:]:  # Skip header line
        # Split on whitespace and extract 11th column (command name)
        parts = line.split(None, 10)  # Split into max 11 parts
        if len(parts) >= 11:
            command = parts[10]
            # Filter by name using Python's 'in' operator instead of shell grep
            if name.lower() in command.lower():
                names.append(command)
    
    return jsonify({"success": True, "names": names})



@bp.route("/deserialized_descr", methods=["POST"])
def deserialized_descr():
    pickled = request.form.get('pickled')
    data = base64.urlsafe_b64decode(pickled)
    # vulnerability: Insecure Deserialization
    deserialized = pickle.loads(data)
    return jsonify({"success": True, "description": str(deserialized)})
