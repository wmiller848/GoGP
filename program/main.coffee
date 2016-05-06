######################
## GoGP Program
######################

DNA = '{{dna}}'

######################
## Stage Spawn
######################
console.log('Setting Argv')
args = process.argv.slice(2)
{{vars}}

console.log('Spawning')
{{spawn}}

######################
## Stage Dieing
######################
process.on('beforeExit', ->
  console.log('Dieing...')
  {{dieing}}
)

######################
## Stage Dead
######################
process.on('exit', ->
  console.log('Dead...')
  {{dead}}
)

######################
## Stage Alive
######################
{{alive}}
