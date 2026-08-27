"""
PyPI Package Publisher for evolution-sdk.
"""

import os
import subprocess
import sys
from pathlib import Path
from dotenv import load_dotenv

if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except Exception:
        pass

env_path = Path(__file__).parent / ".env"
load_dotenv(env_path)

token = os.environ.get("PYPI_API_TOKEN")
if not token:
    print("❌ PYPI_API_TOKEN not found in playground/.env.")
    sys.exit(1)

dist_dir = Path(__file__).parent.parent / "sdk" / "python" / "dist"
files = list(dist_dir.glob("evolution_sdk-0.8.0*"))

if not files:
    print(f"❌ No distribution files found in {dist_dir}")
    sys.exit(1)

print(f"📦 Found {len(files)} distribution packages:")
for f in files:
    print(f"  - {f.name}")

# Prepare environment with TWINE credentials
env = os.environ.copy()
env["TWINE_USERNAME"] = "__token__"
env["TWINE_PASSWORD"] = token.strip()

print("\n🚀 Uploading evolution-sdk v0.8.0 to PyPI...")
cmd = [
    sys.executable,
    "-m",
    "twine",
    "upload",
    *[str(f) for f in files],
    "--non-interactive",
]

result = subprocess.run(cmd, env=env)
if result.returncode == 0:
    print("\n🎉 Successfully published evolution-sdk v0.8.0 to PyPI!")
    print("🌐 View on PyPI: https://pypi.org/project/evolution-sdk/")
else:
    print(f"\n❌ PyPI upload failed with exit code {result.returncode}")
    sys.exit(result.returncode)
