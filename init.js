// admin
db.auth("eobadm", "eobpass")

// user
userdb = db.getSiblingDB("eob_system")
userdb.createUser({
  "user": "eobuser",
  "pwd" : "eobuserpass",
  "roles": [
    { "role" : "readWrite", "db" : "eob_system"}
  ],
  "mechanisms": [ "SCRAM-SHA-1" ],
  "passwordDigestor": "client"
})
userdb.auth("eobuser", "eobuserpass")