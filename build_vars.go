package main

// must set these with build options
var BUILD_COMMAND_B64 string
var BUILD_MODE string

var APPNAME string
var HOME string
var PORT string
var COOKIE_EXPIRY_HOURS string
var COOKIE_IDLE_TIMEOUT_HOURS string
var CONFDIR string
var DATADIR string

var CONFIG_FILE = CONFDIR + "/webapp.cfg"
