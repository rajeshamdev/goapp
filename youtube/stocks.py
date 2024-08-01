'''
Tested with python 3.11.8
pip3 install pandas numpy yfinance tabulate
'''

import yfinance as yf
import pandas as pd
from tabulate import tabulate

class StockAnalysis:

    def __init__(self, ticker):
        self.ticker = ticker
        self.df = pd.DataFrame()

    def convertValue(self, value):
        if value >= 1e9:
            return f'{value / 1e9:.2f}B'  # Convert to billions
        elif value >= 1e6:
            return f'{value / 1e6:.2f}M'  # Convert to millions
        elif value >= 1e3:
            return f'{value / 1e3:.2f}K'  # Convert to millions
        else:
            return f'{value:0.3f}'  # Return as is with 3 decimal places

    def getFundamentalMetrics(self):
        '''
        Get Fundamental metrics of stock. These include:
            - Revenue
            - Profit
            - num of shares
            - FCF (free cash flow)
            - cash and equivalent
            - Loan to Debt Ratio
        Args:
            ticker (str): stock ticker per New York stock exchange

        Returns:
            pandas df: returns pandas data frame
        '''

        stock = yf.Ticker(self.ticker)
        quarters = 4

        '''
        Step-1:
        From Balance Sheet, find:
        - Loan to Debt Ratio
        - Cash And Equivalwent
        '''
        balanceSheet = stock.quarterly_balance_sheet.T.head(quarters)
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

        '''
        Step-2: From financial and cash flow statement find the remaining metrics.
        '''
        financials = stock.quarterly_financials.T.head(quarters)
        cashFlow = stock.quarterly_cash_flow.T.head(quarters)
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

    def fetchStockData(self, startDate, endDate):
        '''
        Populate the data frame.

        Args:
            startDate (str): Start Date in the format of YYYY-MM-DD
            EndDate (str): End Date in the format of YYYY-MM-DD

        Returns:
            None
        '''
        self.df = yf.download(self.ticker, start=startDate, end=endDate)

    def findRSI(self, period=14):
        '''
        Calculate RSI .

        Args:
            period (int): Time frame for RSI
        Returns:
            None
        '''
        delta = self.df['Close'].diff()
        gain = (delta.where(delta > 0, 0)).rolling(window=period).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(window=period).mean()
        rs = gain / loss
        rsi = 100 - (100 / (1 + rs))
        self.df['RSI'] = rsi

    def findMACD(self, shortWindow=12, longWindow=26, signalWindow=9):
        '''
        Calculate MACD

        Args:
            period (int): Time frame for RSI
        Returns:
            None
        '''
        self.df['EMA_short'] = self.df['Close'].ewm(span=shortWindow, adjust=False).mean()
        self.df['EMA_long'] = self.df['Close'].ewm(span=longWindow, adjust=False).mean()
        self.df['MACD'] = self.df['EMA_short'] - self.df['EMA_long']
        self.df['Signal'] = self.df['MACD'].ewm(span=signalWindow, adjust=False).mean()
        self.df['MACD_Histogram'] = self.df['MACD'] - self.df['Signal']


def main():

    ticker = input("Enter Stock Ticker:  ")
    stock = StockAnalysis(ticker)

    # get fundamental metrics
    metrics = stock.getFundamentalMetrics()
    # print(tabulate(Metrics, headers='keys', tablefmt='grid'))
    metrics = metrics.map(stock.convertValue)
    print(metrics)

    # find RSI, MACD etc
    startDate = input("Enter Start Date in YYYY-MM-DD format: ")
    endDate = input("Enter End Date in YYYY-MM-DD format: ")
    stock.fetchStockData(startDate, endDate)
    stock.findRSI()
    stock.findMACD()
    # Display the last few rows with the indicators
    print(stock.df[['Close', 'RSI', 'MACD', 'Signal', 'MACD_Histogram']].tail(20))


if __name__ == "__main__":
    main()
