import yfinance as yf
import pandas as pd
from tabulate import tabulate


def convert_value(value):
    if value >= 1e9:
        return f'{value / 1e9:.2f}B'  # Convert to billions
    elif value >= 1e6:
        return f'{value / 1e6:.2f}M'  # Convert to millions
    elif value >= 1e3:
        return f'{value / 1e3:.2f}K'  # Convert to millions
    else:
        return f'{value:0.3f}'  # Return as is with 3 decimal places


def get_financial_data(ticker):

    stock = yf.Ticker(ticker)

    '''
    From Balance Sheet, find:
      - Loan to Debt Ratio
      - Cash And Equivalwent
    '''
    balanceSheet = stock.quarterly_balance_sheet.T.head(4)
    cashAndEq = balanceSheet.get('Cash And Cash Equivalents', None)
    LTDRatio = []
    # find Loan To Debt Ratio
    for index, row in balanceSheet.iterrows():
        totalDebt = row.get('Total Liabilities Net Minority Interest', None)
        longTermDebt = row.get('Long Term Debt', None)  # Assuming long-term debt includes loans

        if totalDebt is None or longTermDebt is None:
            LTDRatio.append("Data not available")
        else:
            # Calculate Loan to Debt Ratio
            if totalDebt == 0:
                LTDRatio.append("Total debt is zero, cannot calculate ratio")
            else:
                loanToDebtRatio = longTermDebt / totalDebt
                LTDRatio.append(loanToDebtRatio)

    financials = stock.quarterly_financials.T.head(4)
    cashFlow = stock.quarterly_cash_flow.T.head(4)
    # Prepare data for output
    data = {
        'Revenue': financials['Total Revenue'].values,
        'Profit': financials['Net Income'].values,
        'shares': financials['Basic Average Shares'].values,
        'EPS': financials['Basic EPS'].values,
        'CashAndEq': cashAndEq,
        'LTDR': LTDRatio,
        'FCF': cashFlow['Free Cash Flow'],
        # 'dil EPS': financials['Diluted EPS'].values,
        # 'EBITDA': financials['EBITDA'].values,
        # 'EBIT': financials['EBIT'].values,
        # 'dil shares' : financials['Diluted Average Shares'].values,
    }

    # Convert to DataFrame for easier readability
    df = pd.DataFrame(data, index=financials.index)
    pd.set_option('display.float_format', '{}'.format)
    return df

def main():
    ticker = input("Enter Stock Ticker:  ")
    Metrics = get_financial_data(ticker)
    # print(tabulate(Metrics, headers='keys', tablefmt='grid'))
    Metrics = Metrics.map(convert_value)
    print(Metrics)

if __name__ == "__main__":
    main()
