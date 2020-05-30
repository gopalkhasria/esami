function register() {
  event.preventDefault();
  //var regexp = /(?=.[A-Z])(?=.[a-z])(?=.[0-9])(?=.[-_+@$&.,:;]).{8,32}/;
  var password = document.getElementById("password").value;
  var email = document.getElementById("email").value;
  var name = document.getElementById("name").value;
  var err = false;
  if (document.getElementById("confirm_password").value != password) {
    document.getElementById("msg").innerHTML = "le password non corrispondono";
    $('#error').show();
    err = true;
  }
  if (password.match < 8) {
    document.getElementById("msg").innerHTML = " Password errata inserire almeno 8 caratteri";
    $('#error').show();
    err = true;
  }
  if (name.length < 5) {
    document.getElementById("msg").innerHTML = " Nome non valido"
    $('#error').show();
    err = true;
  }
  var mailformat = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
  if (email.match(mailformat)) { }
  else {
    document.getElementById("msg").innerHTML = " Email non valido"
    $('#error').show();
    err = true;
  }
  if (!err) {
    document.getElementById("register").submit();
  }
}
