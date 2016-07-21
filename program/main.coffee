#!{{coffee_path}}

######################
## GoGP Program
## Version - 0.1.0
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

assertMap =
{{assertMap}}
##
##
match = (output) ->
  diff = Number.MAX_VALUE
  ok = ""
  ov = 0
  for k, v of assertMap
    d = Math.abs(v - output)
    if d < diff
      diff = d
      ok = k
      ov = v
  # "#{ok} (#{ov}) from #{output}"
  "#{ok}"
##
##
run = ->
  {{vars}}
  output = {{output}}
  if isNaN(output)
    output = ''
  ot = match(output)
  process.stdout.write(new Buffer.from(ot.toString() + '\n'))

if pargs.length == 0
  process.stdin.setEncoding('utf8')
  process.stdin.on('readable', ->
    chunks = process.stdin.read()
    if chunks
      chunks = chunks.toString().trim().split('\n')
      for chunk in chunks
        args = chunk.split(',')
        run()
  )
  process.stdin.on('end', ->
    #process.stdout.write(new Buffer.from('\r\n'))
  )
else
  args = pargs.join(' ').trim().split(' ')
  run()
