
import { firebaseAuth } from "../../lib/authFirebase"

firebaseAuth.auth().onAuthStateChanged(async (user) => {
  if (
    !user &&
    !window.location.pathname.startsWith("/auth/login") &&
    !window.location.pathname.startsWith("/auth/criar")
  ) {
    console.log("user is null")
    localStorage.removeItem("token")
    window.location.href = "/auth/login"
    return
  }
  if (
    user &&
    (window.location.pathname.startsWith("/auth/login") ||
      window.location.pathname.startsWith("/auth/criar"))
  ) {
    window.location.href = "/"
  }
  await user
    ?.getIdToken()
    .then((token) => {
      localStorage.setItem("token", token)
    })
    .catch((error) => {
      console.log(error)
    })

  console.log("onAuthStateChanged: ", user?.displayName)
})

