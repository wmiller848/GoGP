#!/usr/local/bin/coffee

######################
## GoGP Program
## Version - 0.1.0
######################

DNA = 'bdc7f223df2fa9194876f426fe389607631b9f0be51c534ed1f9ab23494d4a5fd1685d1668c2863dd837f9cd87da9199591a28242342aa6347411bd1251a5ee6770dece01a4f7071626f8774ce8ac81f0449ca3f41b4d95298c8c9098b014674ca4340d11dd9a18a9769803bf1d8a330d63da90c86b0|d3826d5e54f7bd28813538d0b684bd38bd8883a3ba870e8a4bca2e4e01c97d12eb97ac5f09361d339a4a52cde46e254f3dbbcf3973c86ef8c05a728b9918e958f03bcfe164092bdfc488d64733e40127db45e4f61b7dddf2b0cbbd8d5ebdd98e5194d382236d5e54f7'
pargs = process.argv.slice(2)
args = null

process.on('beforeExit', ->
  #console.log('Dieing...')
)

process.on('exit', ->
  #console.log('Dead...')
)

inputMap =
  'noInputStrings': null

assertMap =
  'Iris-virginica': 24
  'Iris-setosa': 8
  'Iris-versicolor': 16

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
  $zz = Number(args[0]);$zy = Number(args[1]);$zx = Number(args[2]);$zw = Number(args[3]);
  $zz = inputMap[args[0]] if isNaN($zz);$zy = inputMap[args[1]] if isNaN($zy);$zx = inputMap[args[2]] if isNaN($zx);$zw = inputMap[args[3]] if isNaN($zw);
  output = 4+(9*$zw)
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

