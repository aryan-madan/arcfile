// Handle file upload
document
  .getElementById("uploadForm")
  .addEventListener("submit", (event) => {
    event.preventDefault();

    const fileInput = document.getElementById("fileInput");
    const emailInput = document.getElementById("emailInput");
    const uploadResult = document.getElementById("uploadResult");

    // Prepare form data
    const formData = new FormData();
    formData.append("file", fileInput.files[0]);
    formData.append("email", emailInput.value);

    // POST to /api/upload
    fetch("/api/upload", {
      method: "POST",
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        // Check for successful upload (code 200)
        if (data.code === 200) {
          uploadResult.innerHTML = `
                          <strong>Success:</strong> ${data.message}<br>
                          <strong>Identifier:</strong> ${data.data.identifier}<br>
                          <strong>Filename:</strong> ${data.data.filename}<br>
                          <strong>Email:</strong> ${data.data.email}
                      `;
        } else {
          // Show error message (for example, missing/invalid email)
          uploadResult.innerHTML = `<strong>Error:</strong> ${data.message}`;
        }
        uploadResult.style.display = "block"; // Show the result
      })
      .catch((error) => {
        uploadResult.innerHTML = `<strong>Error:</strong> ${error}`;
        uploadResult.style.display = "block"; // Show the error
      });
  });

// Handle file download
document
  .getElementById("downloadForm")
  .addEventListener("submit", (event) => {
    event.preventDefault();

    const identifier = document.getElementById("identifierDownload").value;
    const downloadResult = document.getElementById("downloadResult");
    const url = `/api/download/${identifier}`;

    fetch(url)
      .then((response) => {
        if (!response.ok) {
          return response.text().then((body) => {
            throw new Error(
              `Failed to download file. Please check the identifier and try again.\nResponse: ${body}`
            );
          });
        }

        // Optionally extract filename from headers (if provided)
        const disposition = response.headers.get("Content-Disposition");
        let filename = "downloaded_file";
        if (disposition && disposition.indexOf("filename=") !== -1) {
          const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
          const matches = filenameRegex.exec(disposition);
          if (matches?.[1]) {
            filename = matches[1].replace(/['"]/g, "");
          }
        }
        return response
          .blob()
          .then((blob) => ({ blob: blob, filename: filename }));
      })
      .then((obj) => {
        // Create a temporary link to trigger download
        const downloadUrl = window.URL.createObjectURL(obj.blob);
        const a = document.createElement("a");
        a.href = downloadUrl;
        a.download = obj.filename;
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(downloadUrl);
        downloadResult.innerHTML = "Download initiated.";
        downloadResult.style.display = "block"; // Show the result
      })
      .catch((error) => {
        downloadResult.innerHTML = `<strong>Error:</strong> ${error.message}`;
        downloadResult.style.display = "block"; // Show the error
      });
  });

document.getElementById("fileInput").addEventListener("change", function () {
  const uploadForm = document.getElementById("uploadForm");
  const uploadText = document.querySelector(".upload-text");
  const uploadArea = document.getElementById("uploadArea");
  const removeFileButton = document.getElementById("removeFileButton");

  if (this.files && this.files.length > 0) {
    uploadForm.style.display = "block";
    uploadText.style.display = "none";
    uploadArea.style.display = "none";
    removeFileButton.style.display = "flex";
  } else {
    uploadForm.style.display = "none";
    uploadText.style.display = "block";
    uploadArea.style.display = "flex";
    removeFileButton.style.display = "none";
  }
});

document
  .getElementById("removeFileButton")
  .addEventListener("click", () => {
    const fileInput = document.getElementById("fileInput");
    const uploadForm = document.getElementById("uploadForm");
    const uploadText = document.querySelector(".upload-text");
    const uploadArea = document.getElementById("uploadArea");
    const removeFileButton = document.getElementById("removeFileButton");

    fileInput.value = "";
    uploadForm.style.display = "none";
    uploadText.style.display = "block";
    uploadArea.style.display = "flex";
    removeFileButton.style.display = "none";
  });
