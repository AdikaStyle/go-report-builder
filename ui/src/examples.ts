

export const ExampleHTML = `
<!DOCTYPE html>
<html>
<head>
	<style>
		.invoice-container {
			padding: 10px;
			display: flex;
			flex-direction: column;
			justify-content: center;
		}

		table, th, td {
			border: 1px solid black;
			border-collapse: collapse;
			text-align: center;
		}

		th.small {
			width: 1%;
		}

		th.large {
			width: 20%;
		}
	</style>
	<title>Example - Invoice</title>
</head>
<body>

<div class="invoice-container" id="printable">
	<h2 style="text-align: center">Invoice # <u>{{ data.invoiceId }}</u></h2>
	<table style="width:100%">
		<thead>
		<tr>
			<th class="large">SKU</th>
			<th class="large">Name</th>
			<th class="small">No. of Units</th>
			<th class="small">Net Weight</th>
			<th class="small">Unit Value</th>
			<th class="small">Total Value</th>
		</tr>
		</thead>

		<tbody>
			{% for l in data.lineItems %}
			<tr>
				<td>{{ l.sku }}</td>
				<td>{{ l.name }}</td>
				<td>{{ l.quantity }}</td>
				<td>{{ l.netWeight }}</td>
				<td>{{ l.unitValue }}</td>
				<td>{{ l.totalValue }}</td>
			</tr>
			{% endfor %}
		</tbody>
		</table>
</div>

</body>
</html>
`

export const ExampleData= `
{
  "invoiceId": "SH200992",
  "lineItems": [
    {
      "name": "Large box of gold",
      "netWeight": "0.1 KG",
      "quantity": 3,
      "sku": "A00005454",
      "totalValue": "30$",
      "unitValue": "10$"
    },
    {
      "name": "Fresh Air",
      "netWeight": "0.3 KG",
      "quantity": 10,
      "sku": "A0000522354",
      "totalValue": "150$",
      "unitValue": "15$"
    }
  ]
}
`