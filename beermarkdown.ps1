param($b)

"---"
"layout: post"
"title: $($b.name)"
"tags: [ beer ]"
"---"

"## Malt"
$b.FERMENTABLES.FERMENTABLE | % { "-  $($_.Name), $($_.DISPLAY_AMOUNT)" }
"## Hops"
$b.Hops.Hop | % { "-  $($_.Name), $($_.DISPLAY_AMOUNT) oz, $($_.TIME) min" }
"## Yeast"
$b.Yeasts.Yeast | % { "-  $($_.Name), $($_.Form)" }
"## Notes"
$b.Notes
