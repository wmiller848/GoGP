######################
## GoGP Program
######################

DNA = '{{dna}}'

######################
## Stage Spawn
######################
#console.log('Setting Argv')
pargs = process.argv.slice(2)
args = null
#console.log('Spawning')
{{spawn}}

######################
## Stage Dieing
######################
process.on('beforeExit', ->
  #console.log('Dieing...')
  {{dieing}}
)

######################
## Stage Dead
######################
process.on('exit', ->
  #console.log('Dead...')
  {{dead}}
)

######################
## Stage Alive
######################
run = ->
  {{vars}}
  output = {{alive}}
  if isNaN(output)
    output = 'NaN'
  process.stdout.write(new Buffer.from(output.toString() + '\n'))

if pargs.length == 0
  process.stdin.setEncoding('utf8')
  process.stdin.on('readable', ->
    chunks = process.stdin.read()
    if chunks
      chunks = chunks.toString().trim().split('\n')
      for chunk in chunks
        args = []
        data = Buffer.from(chunk.split(' '))
        args.push(data...)
        run()
  )
  process.stdin.on('end', ->
    #process.stdout.write(new Buffer.from('\r\n'))
  )
else
  args = []
  data = Buffer.from(pargs.join(' ').trim().split(' '))
  args.push(data...)
  run()
