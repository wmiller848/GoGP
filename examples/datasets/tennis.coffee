#!/usr/local/bin/coffee

######################
## GoGP Program
## Version - 0.1.0
######################

DNA = '4fe887be8474977a9f066fab66935775da38f2a280e5da15a246f4bfa0158a03408a59609ede8b4c0bcf7ce15d6c1bb07b2fb0a84f856d870b1586ce84246ffd6a2bb575a831f859fc551f80a4fa9d1ba2d5463a7bb7eb8a8b532e472c1726e0fb66d81585d60a37533293d337319767e578fc48f3800c25a71a8326b7a2463dcb4e8aeb8ac0b2404232221bef4f7da7be977ad82babb893d9cf8353f6da38cdb3a2809a049746cb8a6503a5135e2a8bfc48947c86fa956b227e1b8744b0bcf32b58906a7430846a2b1ab27f75577870da2c4ccd704cdae5c716c8ef872f1bdf46cb7bb7eb8a8bbfe62e8cae4385e21585247a4475fd5731e548b3063cbe636ffeab66935731174338b3a280e59a26acb704fda24641116e34cd8a03668a5b4268221b87057bb0bc96858087ce846f2bb25731dae67ccd4efc354f4fe887be74977a9faa6fab669357e4da2538b3a280f6e59ada15a246f4bfa0158af203409a59609e8bd60bcf7ce15d236c1bb07b2fb0a84f856a870b1586ce84246fed6a2bb5f075a87f31f859fc551f8009fa9d1ba2d5463a7bb71deb8a8b532e472c1726e0fbd81585d67a37533293d3373167e5781f48b3fc80e50cf1a71a260fb789a21a463dcb4e218aeb8aa4b2404232221b4f68be9730d82babb893cf835331da38cdb3a2809a19049746cbd68a6503a513fb5e8a8bfc48|d370c53db6f20a312b3078a0c5f47e395154750d55575263191c45e1a5d409106fdd03a2d10de93ab19966709bad4e3d5aa60a7787310b30d38bc53db6f20a31e12b3078a0c5f41ac854750d55575263196845e1a5d4096fdd03cca2d1e0e9613ab199'
pargs = process.argv.slice(2)
args = null

process.on('beforeExit', ->
  #console.log('Dieing...')
)

process.on('exit', ->
  #console.log('Dead...')
)

inputMap =
  'sunny': 8
  'hot': 8
  'true': 16
  'rainy': 24
  'mild': 16
  'high': 8
  'false': 8
  'overcast': 16
  'cool': 24
  'normal': 16

assertMap =
  'no': 8
  'yes': 16

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
  output = 13^($zz&971&($zx|($zy^510^31^$zw^7^($zz^($zx|$zz)))))
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

