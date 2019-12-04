#Download Source
spectool -g tism.spec

#Extract Source
fedpkg --release f25 prep

#Build
fedpkg --release f25 local
