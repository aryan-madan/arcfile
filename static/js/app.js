const form = document.getElementById("form");
const inputFile = document.getElementById("file");
const emailInput = document.getElementById("email"); // Add this reference

const handleSubmit = (event) => {
  event.preventDefault(); // Critical: prevent default form submission

  const formData = new FormData();

  // Append file(s)
  for (const file of inputFile.files) {
    formData.append("file", file); // Key must match server's "file"
  }

  // Append email
  formData.append("email", emailInput.value);

  // Send request
  fetch("/api/upload", {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => console.log(data))
    .catch((error) => console.error("Error:", error));
};

form.addEventListener("submit", handleSubmit);
