import 'bootstrap/dist/css/bootstrap.min.css'
import './index.css';
import '@fortawesome/fontawesome-free/css/fontawesome.min.css'
import '@fortawesome/fontawesome-free/css/solid.min.css'
import { App } from './app';
import jQuery from "jquery";

// Force declarationn of jquery
Object.assign(window, { $: jQuery, jQuery })

window.addEventListener("load", function() {
  App.start();
});