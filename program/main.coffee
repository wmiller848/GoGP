#!{{coffee_path}}

######################
## GoGP Program
## Version - 0.0.1
######################

DNA = '{{dna}}'
pargs = process.argv.slice(2)
args = null

process.on('beforeExit', ->
  #console.log('Dieing...')
)

process.on('exit', ->
  #console.log('Dead...')
)

##
##
run = ->
  {{vars}}
  output = {{output}}
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
        data = new Buffer.from(chunk.split(' '))
        args.push(data...)
        run()
  )
  process.stdin.on('end', ->
    #process.stdout.write(new Buffer.from('\r\n'))
  )
else
  args = []
  data = new Buffer.from(pargs.join(' ').trim().split(' '))
  args.push(data...)
  run()
