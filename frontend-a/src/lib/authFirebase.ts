
import 'firebase/compat/auth'
import firebase from 'firebase/compat/app'

// TODO: Replace the following with your app's Firebase project configuration
const firebaseConfig = {
    apiKey: 'AIzaSyD9zukt9KOqE9rLQMdw6GripeqfMIy0PpI',
    authDomain: 'sos-do-maceneiro.firebaseapp.com',
    projectId: 'sos-do-maceneiro',
    storageBucket: 'sos-do-maceneiro.appspot.com',
    messagingSenderId: '814307894175',
    appId: '1:814307894175:web:775223d738d69750a693aa',
    measurementId: 'G-QNJVSJPCWY'
}

export const firebaseAuth = firebase.initializeApp(firebaseConfig)
