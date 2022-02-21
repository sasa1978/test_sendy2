async function submit() {
  setLoading(true);
  const url = document.getElementById("list-id").value;
  if (!url) {
    error("Playlist ID не может быть пустым!")
    return;
  }
  try {
    const resp = await fetch('/getPlayList', {
      method: "POST",
      body: url
    });
    if (!resp.ok) {
      error("Не могу получить данные :'(")
      return;
    }
    const txtImgRaw = await resp.text();
    console.log(txtImgRaw.split("\n").length);
    const txtImgCleaned = txtImgRaw.replaceAll("\n", "<br/>");
    console.log(txtImgCleaned.split("<br/>").length);
    console.log(txtImgCleaned)
    document.getElementById("txtimg").innerHTML = txtImgCleaned;
    succ("Список успешно получен!");
  } catch (e) {
    error("Не могу получить данные :'(")
  }
}

function clearStatus() {
  document.getElementById("alert").innerHTML = "";
}

function status(msg, type) {
  document.getElementById("alert").innerHTML = `
    <div class="alert alert-${type}" role="alert">
      ${msg}
    </div>
  `;
}

function warn(msg) {
  status(msg, "warning");
}

function succ(msg) {
  status(msg, "success");
}

function info(msg) {
  status(msg, "info");
}

function error(msg) {
  status(msg, "danger");
}

function setLoading(isLoading) {
  if (isLoading) {
    const spinner = `
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>`;
    status(spinner, "info");
  } else {
    clearStatus();
  }
}