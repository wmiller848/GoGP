main = ->
  console.log('Running main')
  args = process.argv.slice(2)
  console.log(args)
  $a = args[0]
  $b = args[1]
  ##
  ## Parse variables from input
  #{{}}
  ## Run the expression
  #{{}}
  $a*(30-10-$b-12-(10*70))

console.log(main())
