#!/usr/local/bin/coffee

######################
## GoGP Program
## Version - 0.1.0
######################

DNA = '463d3d0ee1bb024ba9ecb52441e7ae39edb592ccb3bd72490d80703e84d396c06628ba649ec05a140b2605ce865b41d02a266ca61dfb1ac87e6a32d6de56349505e6c039c62fd5d65d6e0e39f2f1650c2c3dd95bd4463f3d0ee1664abb4bacb8ec24e7ae39d8edabb592acbdb10d806b3e9b96ba2829ba0d27c05a140b1f26cee305ce6f865b41ee67266c46593d0ee1a646db024ba9ec442441e7ae1c39edb592ccb3bd72b10d803e84d3966628ba649ec0935a140b26ce61ce865b41672a266ca61df2fb1ac86a3256349505e6c039c62fd50ad65d6e0e39f2f1c20c5d2c49085bd4463f3de166|fbafd1c0253103dad61c3fe45b43a633767fab866d8875ae7c0a992a2d9bf8a0b42b69912ed29a8c87ad8f361d08770e3f6b12835209eaa93a25fc5901b5a73638b6371fd4a5f3674119e029fcd597d2681f7ecde3cc7687f2'
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
  'republican': 8
  'democrat': 16

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
  $zz = Number(args[0]);$zy = Number(args[1]);$zx = Number(args[2]);$zw = Number(args[3]);$zv = Number(args[4]);$zu = Number(args[5]);$zt = Number(args[6]);$zs = Number(args[7]);$zr = Number(args[8]);$zq = Number(args[9]);$za = Number(args[10]);$zb = Number(args[11]);$zc = Number(args[12]);$zd = Number(args[13]);$ze = Number(args[14]);$zf = Number(args[15]);
  $zz = inputMap[args[0]] if isNaN($zz);$zy = inputMap[args[1]] if isNaN($zy);$zx = inputMap[args[2]] if isNaN($zx);$zw = inputMap[args[3]] if isNaN($zw);$zv = inputMap[args[4]] if isNaN($zv);$zu = inputMap[args[5]] if isNaN($zu);$zt = inputMap[args[6]] if isNaN($zt);$zs = inputMap[args[7]] if isNaN($zs);$zr = inputMap[args[8]] if isNaN($zr);$zq = inputMap[args[9]] if isNaN($zq);$za = inputMap[args[10]] if isNaN($za);$zb = inputMap[args[11]] if isNaN($zb);$zc = inputMap[args[12]] if isNaN($zc);$zd = inputMap[args[13]] if isNaN($zd);$ze = inputMap[args[14]] if isNaN($ze);$zf = inputMap[args[15]] if isNaN($zf);
  output = 7+($zq+(8&$zv))
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

