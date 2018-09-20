const file = (document.getElementById("file")! as HTMLInputElement)!;
file.addEventListener('change', () => {
  if(file.files!.length > 0) {
    document.getElementById('file-name')!.innerHTML = file.files![0]!.name;
  }
}, false);
